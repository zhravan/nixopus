---
layout: home

hero:
  name: "Nixopus Documentation"
  text: ""
  tagline: All the information you need to know about Nixopus
---
<div class="command-block-container">
  <div class="command-block">
    <div class="command-text">
      curl -sSL https://install.nixopus.com | bash
    </div>
    <button class="copy-button" onclick="navigator.clipboard.writeText(this.previousElementSibling.textContent.trim()); this.innerText='Copied!'; setTimeout(()=>this.innerText='Copy',1200);">Copy</button>
  </div>
</div>
<SponsorsMarquee />
