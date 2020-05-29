/**
 * Utils
 */

// Add code-copy buttons using progressive enhancement
// © 2020. Pranjal Agrawal
// https://www.fiznool.com/blog/2018/09/14/adding-click-to-copy-buttons-to-a-hugo-powered-blog/
(function() {
  'use strict';

  function addCodeControl(containerEl) {
    if (containerEl.querySelector("pre code") == null) {
      return;
    }
    var copyDiv = document.createElement("div");
    copyDiv.className = "code-control";
    copyDiv.innerHTML = "<i data-tippy-content=\"Copy code to clipboard\" class=\"code-copy fa fa-copy\"></i>";
    containerEl.children[0].appendChild(copyDiv);
  }

  // Add copy button to code blocks
  var highlightBlocks = document.getElementsByClassName('highlight');
  Array.prototype.forEach.call(highlightBlocks, addCodeControl);

  // Add clipboard JS controls
  var clipboard = new ClipboardJS('.code-control', {
    target: function(trigger) {
        return trigger.previousElementSibling;
    }
  });

  clipboard.on('success', function(e) {
      var tippyTip = tippy(e.trigger, {content: 'Copied!', trigger: 'manual'});
      tippyTip.show();
      setTimeout(() => tippyTip.destroy(), 2000);
      e.clearSelection();
  });

  clipboard.on('error', function(e) {
      console.error('Action:', e.action);
      console.error('Trigger:', e.trigger);
  });

})();