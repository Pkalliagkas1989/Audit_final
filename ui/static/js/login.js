let API_CONFIG;

async function loadConfig() {
  const res = await fetch("/config");
  if (!res.ok) throw new Error("Failed to load config");
  API_CONFIG = await res.json();
}

window.addEventListener("pageshow", function (event) {
  const form = document.getElementById("loginForm");
  if (form) {
    form.reset();
    const message = document.getElementById("message");
    if (message) {
      message.textContent = "";
    }
  }
});

document
  .getElementById("loginForm")
  .addEventListener("submit", async function (e) {
    e.preventDefault();

    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;
    const message = document.getElementById("message");

    message.textContent = "";
    message.style.color = "red";

    if (!email || !password) {
      message.textContent = "Email and password are required!";
      return;
    }

    try {
      await loadConfig();
      const res = await fetch(API_CONFIG.LoginURI, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ email, password }),
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(`Error: ${res.status}, ${errorText}`);
      }

      const result = await res.json();
      message.style.color = "green";
      message.textContent = "Login successful!";
      sessionStorage.setItem("csrf_token", result.csrf_token);
      localStorage.setItem(
        "session_status",
        JSON.stringify({ status: "logged_in", timestamp: Date.now() })
      );

      window.location.replace("/user");
    } catch (err) {
      console.error(err);
      message.textContent = err.message;
    }
  });
