let API_CONFIG;

async function loadConfig() {
  const res = await fetch("/config");
  if (!res.ok) throw new Error("Failed to load config");
  API_CONFIG = await res.json();
}

document
  .getElementById("registerForm")
  .addEventListener("submit", async function (e) {
    e.preventDefault();

    const username = document.getElementById("username").value.trim();
    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value;
    const confirmPassword = document.getElementById("confirmPassword").value;
    const message = document.getElementById("message");

    message.textContent = "";
    message.style.color = "red";

    if (password !== confirmPassword) {
      message.textContent = "Passwords do not match!";
      return;
    }

    try {
      await loadConfig();
      const res = await fetch(API_CONFIG.RegisterURI, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, email, password }),
      });

      if (!res.ok) {
        const errorText = await res.json();
        throw new Error(`Error: ${res.status}, ${errorText.message}`);
      }

      message.style.color = "green";
      message.textContent = "Registration successful!";
    } catch (err) {
      console.error("Error during registration:", err);
      message.textContent = err.message;
    }
  });
