// Comment handling functionality
class CommentHandler {
  constructor(configManager, dataManager) {
    this.configManager = configManager;
    this.dataManager = dataManager;
  }

  addCommentInput(postContainer, post, refreshFn) {
    const commentBtn = postContainer.querySelector(".comment-btn");
    commentBtn.addEventListener("click", () => {
      let commentInput = postContainer.querySelector(".comment-input");
      if (!commentInput) {
        commentInput = document.createElement("div");
        commentInput.classList.add("comment-input");
        const textarea = document.createElement("textarea");
        textarea.placeholder = "Write your comment...";
        textarea.rows = 3;
        textarea.maxLength = 200;
        textarea.classList.add("comment-input");

        const counter = document.createElement("div");
        counter.style.textAlign = "right";
        counter.style.fontSize = "0.8em";
        counter.style.color = "var(--text-muted)";
        counter.style.marginTop = "0.3em";

        const MAX_LENGTH = 200;
        textarea.addEventListener("input", () => {
          const len = textarea.value.length;
          counter.textContent = `${len} / ${MAX_LENGTH}`;
          if (len > MAX_LENGTH) {
            textarea.style.border = "1px solid var(--color-warning)";
            counter.style.color = "var(--color-warning)";
          } else {
            textarea.style.border = "1px solid var(--bg-quaternary)";
            counter.style.color = "var(--text-muted)";
          }
        });

        const submitBtn = document.createElement("button");
        submitBtn.textContent = "Submit Comment";
        submitBtn.classList.add("comment-btn");
        submitBtn.addEventListener("click", async () => {
          const content = textarea.value.trim();
          let warning = commentInput.querySelector(".comment-warning");
          if (!warning) {
            warning = document.createElement("div");
            warning.className = "comment-warning";
            warning.style.color = "red";
            warning.style.marginTop = "5px";
            commentInput.appendChild(warning);
          }
          warning.textContent = "";
          if (!content) {
            warning.textContent = "Comment cannot be empty.";
            return;
          }
          try {
            await this.configManager.loadConfig();
            const API_CONFIG = this.configManager.getConfig();

            const res = await fetch(API_CONFIG.CommentsURI, {
              method: "POST",
              credentials: "include",
              headers: {
                "Content-Type": "application/json",
                "X-CSRF-Token": sessionStorage.getItem("csrf_token"),
              },
              body: JSON.stringify({
                post_id: post.id,
                content: content,
              }),
            });
            if (!res.ok) {
              const errorText = await res.json();
              throw new Error(
                `Server error: ${res.status}, message: ${errorText.message}`
              );
            }
            const newComment = await res.json();
            textarea.value = "";

            this.dataManager.addCommentToPost(post.id, newComment);
            refreshFn();
          } catch (err) {
            console.error("Failed to submit comment:", err);
          }
        });
        commentInput.appendChild(textarea);
        commentInput.appendChild(counter);
        commentInput.appendChild(submitBtn);
        postContainer.appendChild(commentInput);
      }
    });
  }
}

export { CommentHandler };
