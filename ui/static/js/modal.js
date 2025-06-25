class Modal {
  constructor(API_CONFIG) {
    this.API_CONFIG = API_CONFIG;
    this.init();
  }

  async init() {
    const createBtn = document.getElementById("create-post-btn");
    const modal = document.getElementById("post-modal");
    const closeBtn = document.querySelector(".close-btn");
    const submitBtn = document.getElementById("submit-post");
    const postTitle = document.getElementById("post-title");
    const postBody = document.getElementById("post-body");

    createBtn.onclick = () => modal.classList.remove("hidden");
    closeBtn.onclick = () => modal.classList.add("hidden");
    submitBtn.onclick = () => this.handleSubmit(postTitle, postBody, modal);

    await this.populateCategories();

    const counter = document.createElement("div");
    counter.style.textAlign = "right";
    counter.style.fontSize = "0.8em";
    counter.style.color = "var(--text-muted)";
    counter.style.marginTop = "0.3em";
    postBody.parentNode.appendChild(counter);

    const MAX_LENGTH = 2000;

    postBody.addEventListener("input", () => {
      const len = postBody.value.length;
      counter.textContent = `${len} / ${MAX_LENGTH}`;

      if (len > MAX_LENGTH) {
        postBody.style.border = "1px solid var(--color-warning)";
        counter.style.color = "var(--color-warning)";
      } else {
        postBody.style.border = "1px solid var(--bg-quaternary)";
        counter.style.color = "var(--text-muted)";
      }
    });
  }

  async populateCategories() {
    try {
      const response = await fetch(this.API_CONFIG.CategoriesURI);

      if (!response.ok) {
        const errorText = await response.json();
        throw new Error(
          `HTTP error! status: ${response.status}, message: ${errorText.message}`
        );
      }

      const data = await response.json();

      const container = document.getElementById("post-category");
      container.innerHTML = "";

      data.forEach((category) => {
        const checkbox = document.createElement("input");
        checkbox.type = "checkbox";
        checkbox.id = `cat-${category.id}`;
        checkbox.name = "categories";
        checkbox.value = category.id;

        const label = document.createElement("label");
        label.htmlFor = checkbox.id;
        label.textContent = category.name;

        const wrapper = document.createElement("div");
        wrapper.classList.add("checkbox-item");
        wrapper.appendChild(checkbox);
        wrapper.appendChild(label);

        container.appendChild(wrapper);
      });
    } catch (error) {
      console.error("Error fetching categories:", error);
    }
  }

  async handleSubmit(postTitle, postBody, modal) {
    const title = postTitle.value.trim();
    const content = postBody.value.trim();
    const checkedCategories = document.querySelectorAll(
      'input[name="categories"]:checked'
    );
    const categoryContainer = document.getElementById("post-category");

    // --- Begin new validation logic for title/content ---
    const titleWarningId = "title-warning";
    const contentWarningId = "content-warning";
    const removeWarning = (id) => {
      const el = document.getElementById(id);
      if (el) el.remove();
    };
    const showWarning = (id, message, inputElement) => {
      if (!document.getElementById(id)) {
        const warning = document.createElement("div");
        warning.id = id;
        warning.textContent = message;
        warning.style.color = "red";
        warning.style.marginTop = "5px";
        inputElement.parentNode.appendChild(warning);
      }
    };

    removeWarning(titleWarningId);
    removeWarning(contentWarningId);

    if (!title) {
      showWarning(titleWarningId, "Title is required", postTitle);
    }

    if (!content) {
      showWarning(contentWarningId, "Details are required", postBody);
    }

    if (!title || !content) return;
    // --- End new validation logic for title/content ---

    if (checkedCategories.length === 0) {
      if (!document.getElementById("category-warning")) {
        const warning = document.createElement("div");
        warning.id = "category-warning";
        warning.textContent = "Please select at least one category";
        warning.style.color = "red";
        warning.style.marginTop = "5px";
        categoryContainer.appendChild(warning);
      }
      return;
    } else {
      const warning = document.getElementById("category-warning");
      if (warning) warning.remove();
    }

    const categoryIds = Array.from(checkedCategories).map((cb) =>
      parseInt(cb.value, 10)
    );
    if (categoryIds.length === 0) return;

    try {
      const response = await fetch(this.API_CONFIG.CreatePostURI, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-CSRF-Token": sessionStorage.getItem("csrf_token"),
        },
        credentials: "include",
        body: JSON.stringify({
          title,
          content,
          category_ids: categoryIds,
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(
          `Server error: ${response.status}, message: ${errorData.message}`
        );
      }

      const result = await response.json();
      console.log("Post created successfully:", result);

      postTitle.value = "";
      postBody.value = "";
      modal.classList.add("hidden");
      window.location.reload();
    } catch (error) {
      console.error("Error creating post:", error);
    }
  }
}

// ✅ Export this
export function initModal(API_CONFIG) {
  new Modal(API_CONFIG);
}
