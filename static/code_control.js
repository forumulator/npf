/**
 * Code blocks control.
 */

function CodeThemeToggler() {
  function addDarkStyles() {
    const link = document.createElement('link');
    link.type = 'text/css';
    link.id = 'dark-theme-styles';
    link.rel = 'stylesheet';
    link.href = '/css/dracula.css';
    link.setAttribute('disabled', '');
    document.head.appendChild(link);
    return link;
  }

  const regular =  document.getElementById('code-styles');
  let dark = document.getElementById('dark-theme-styles');
  if (!dark) {
    dark = addDarkStyles();
  }

  return {
    toggle: function() {
      if (dark.hasAttribute('disabled')) {
        dark.removeAttribute('disabled');
        regular.setAttribute('disabled', '');
      } else {
        dark.setAttribute('disabled', 'true');
        regular.removeAttribute('disabled');
      }
    }
  }
}

// Add code-copy buttons using Clipboard JS.
(function () {
  'use strict';

  const CODE_COPY = '<span data-tippy-content="Copy code to clipboard" class="code-copy material-icons">'
      + 'content_copy</span>';

  const DARK_THEME = '<span data-tippy-content="Toggle code theme" class="dark-theme material-icons">'
      + 'brightness_medium</span>';

  function addCodeControl(containerEl) {
    if (containerEl.querySelector("pre code") == null) {
      return;
    }
    var copyDiv = document.createElement("div");
    copyDiv.className = "code-control";
    // copyDiv.innerHTML = "<i  class=\"code-copy fa fa-copy\"></i>";
    copyDiv.innerHTML = DARK_THEME + CODE_COPY;
    containerEl.children[0].appendChild(copyDiv);
  }

  // Add copy button to code blocks
  var highlightBlocks = document.getElementsByClassName('highlight');
  Array.prototype.forEach.call(highlightBlocks, addCodeControl);

  // Add clipboard JS controls
  var clipboard = new ClipboardJS('.code-copy', {
    target: function (trigger) {
      return trigger.parentElement.previousElementSibling;
    }
  });

  clipboard.on('success', function (e) {
    let tippyTip = tippy(e.trigger.parentElement, {content: 'Copied!', trigger: 'manual'});
    tippyTip.show();
    setTimeout(() => tippyTip.destroy(), 2000);
    e.clearSelection();
  });

  clipboard.on('error', function (e) {
    console.error('Action:', e.action);
    console.error('Trigger:', e.trigger);
  });

  // Add theme toggle handler.
  document.querySelectorAll('.code-control .dark-theme')
      .forEach(item => {
        item.addEventListener('click', () => { CodeThemeToggler().toggle(); });
      })
})();
