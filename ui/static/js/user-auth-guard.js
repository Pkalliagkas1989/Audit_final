// Authentication guard functionality
class AuthGuard {
  constructor() {
    this.hideBodyImmediately();
  }

  hideBodyImmediately() {
    // Hide body immediately to prevent flash of protected content (with unique id)
    const style = document.createElement("style");
    style.id = "auth-guard-style";
    style.innerHTML = "body { display: none !important; }";
    document.head.appendChild(style);
  }

  async checkAuthentication(API_CONFIG) {
    try {
      const res = await fetch(API_CONFIG.AuthURI, {
        credentials: "include",
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(`Error: ${res.status}, ${errorText}`);
      }

      const sessionData = await res.json();
      this.session = sessionData;

      const guardStyle = document.getElementById("auth-guard-style");
      if (guardStyle) guardStyle.remove();
      document.body.style.display = "";

      return true;
    } catch (err) {
      window.location.href = "/";
      console.error("Authentication failed:", err.message);
      return false;
    }
  }
}

export { AuthGuard };
