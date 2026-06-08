<script>
  import { marked } from 'marked';
  import hljs from 'highlight.js';
  import 'highlight.js/styles/github-dark.css';
  import { onMount } from 'svelte';

  export let content = '';

  $: html = marked.parse(content);

  function copyToClipboard(text, btn) {
    navigator.clipboard.writeText(text).then(() => {
      const originalText = btn.innerText;
      btn.innerText = 'Copied!';
      setTimeout(() => {
        btn.innerText = originalText;
      }, 2000);
    });
  }

  onMount(() => {
    // Initial highlighting
    document.querySelectorAll('pre code').forEach((block) => {
      hljs.highlightElement(block);
    });
  });

  // Re-run highlighting when content changes
  $: if (content) {
    setTimeout(() => {
      document.querySelectorAll('pre code').forEach((block) => {
        if (!block.dataset.highlighted) {
          hljs.highlightElement(block);
          block.dataset.highlighted = 'true';

          // Add copy button if not exists
          const pre = block.parentElement;
          if (pre && !pre.querySelector('.copy-btn')) {
            const btn = document.createElement('button');
            btn.className = 'copy-btn absolute top-2 right-2 px-2 py-1 text-xs bg-gray-700 hover:bg-gray-600 rounded text-gray-300 transition-colors';
            btn.innerText = 'Copy';
            btn.onclick = () => copyToClipboard(block.innerText, btn);
            pre.appendChild(btn);
          }
        }
      });
    }, 0);
  }
</script>

<div class="prose prose-invert max-w-none">
  {@html html}
</div>

<style>
  :global(.prose pre) {
    position: relative;
    background-color: #1e1e1e !important;
  }

  :global(.prose p) {
    margin-bottom: 1rem;
    line-height: 1.6;
  }

  :global(.prose ul) {
    list-style-type: disc;
    padding-left: 1.5rem;
    margin-bottom: 1rem;
  }
</style>
