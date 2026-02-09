const pages = [
  {
    title: "Open Source Projects Hub",
    url: "https://github.com/topics/open-source",
    snippet:
      "Explore popular open-source projects, trending repositories, and collaboration tools.",
  },
  {
    title: "MDN Web Docs",
    url: "https://developer.mozilla.org/",
    snippet:
      "Authoritative resources for HTML, CSS, JavaScript, and modern web platform APIs.",
  },
  {
    title: "Tech News Today",
    url: "https://news.ycombinator.com/",
    snippet:
      "Read and discuss the latest technology, startup, and programming news.",
  },
  {
    title: "Learn JavaScript",
    url: "https://javascript.info/",
    snippet:
      "A modern tutorial from basics to advanced topics with practical examples.",
  },
  {
    title: "Design Inspiration Gallery",
    url: "https://dribbble.com/",
    snippet:
      "Discover UI and product design inspiration from global creative communities.",
  },
  {
    title: "Web Performance Guide",
    url: "https://web.dev/",
    snippet:
      "Guidance and best practices for building fast, reliable, and engaging websites.",
  },
];

function scoreResult(query, page) {
  const words = query.toLowerCase().split(/\s+/).filter(Boolean);
  const haystack = `${page.title} ${page.snippet}`.toLowerCase();

  return words.reduce((score, word) => {
    if (page.title.toLowerCase().includes(word)) return score + 5;
    if (haystack.includes(word)) return score + 2;
    return score;
  }, 0);
}

function findResults(query) {
  return pages
    .map((page) => ({ ...page, score: scoreResult(query, page) }))
    .filter((page) => page.score > 0)
    .sort((a, b) => b.score - a.score)
    .slice(0, 6);
}

function boot() {
  const searchForm = document.querySelector("#search-form");
  const queryInput = document.querySelector("#query");
  const resultsRoot = document.querySelector("#results");
  const luckyBtn = document.querySelector("#lucky-btn");
  const searchBtn = document.querySelector("#search-btn");
  const resultTemplate = document.querySelector("#result-template");

  if (!searchForm || !queryInput || !resultsRoot || !luckyBtn || !searchBtn || !resultTemplate) {
    return;
  }

  function renderStatus(message) {
    resultsRoot.innerHTML = `<p class="status">${message}</p>`;
  }

  function renderResults(query) {
    const results = findResults(query);

    if (!results.length) {
      renderStatus(`No results found for "${query}".`);
      return;
    }

    resultsRoot.innerHTML = "";

    for (const result of results) {
      const node = resultTemplate.content.cloneNode(true);
      const [urlLink, titleLink] = node.querySelectorAll("a");

      urlLink.href = result.url;
      urlLink.textContent = result.url;
      titleLink.href = result.url;
      titleLink.textContent = result.title;
      node.querySelector(".result-snippet").textContent = result.snippet;

      resultsRoot.appendChild(node);
    }
  }

  function runSearch() {
    const query = queryInput.value.trim();

    if (!query) {
      renderStatus("Type a search term to begin.");
      return;
    }

    renderResults(query);
  }

  searchForm.addEventListener("submit", (event) => {
    event.preventDefault();
    runSearch();
  });

  searchBtn.addEventListener("click", (event) => {
    event.preventDefault();
    runSearch();
  });

  luckyBtn.addEventListener("click", () => {
    const query = queryInput.value.trim();

    if (!query) {
      renderStatus("Type a search term before trying luck.");
      return;
    }

    const [first] = findResults(query);

    if (!first) {
      renderStatus(`No lucky result for "${query}".`);
      return;
    }

    window.location.href = first.url;
  });

  renderStatus("Search for topics like 'javascript', 'design', or 'open source'.");
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", boot);
} else {
  boot();
}
