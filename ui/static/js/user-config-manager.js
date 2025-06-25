// Configuration management
class ConfigManager {
  constructor() {
    this.API_CONFIG = null;
  }

  async loadConfig() {
    const res = await fetch("/config");
    if (!res.ok) throw new Error("Failed to load config");
    this.API_CONFIG = await res.json();
    return this.API_CONFIG;
  }

  getConfig() {
    return this.API_CONFIG;
  }
}

export { ConfigManager };