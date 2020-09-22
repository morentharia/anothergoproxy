(function () {
  const ANOTHERPROXY_API_URL = "{{{ANOTHERPROXY_API_URL}}}";

  function ProxyLog(t, p) {
    let dataToSend = { type: t, params: p };
    fetch(ANOTHERPROXY_API_URL + "/log", {
      method: "post",
      // mode: "no-cors",
      mode: "cors",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(dataToSend),
    })
      .then((response) => {
        if (response.status !== 200) {
          // log("Status Code: " + response.status);
          return result;
        }
        //
        // // Examine the text in the response
        // response.json().then(function (data) {
        //   // log(data);
        // });
      })
      .catch((err) => {
        // log("errrorrrr" + err);
      });
  }

  if (typeof window.ANOTHERPROXY_FLAG === "undefined") {
    var wraplogfunc = (log, ...arg) => {
      let result = log(...arg);
      ProxyLog("console.log", arg);
      return result;
    };
    console.log = _.wrap(console.log, wraplogfunc);
    window.ANOTHERPROXY_FLAG = true;

    document.addEventListener("DOMContentLoaded", function (event) {
      var mutationObserver = new MutationObserver(function (mutations) {
        mutations.forEach(function (mutation) {
          // if (mutation.type === "childList") {
          //   // console.log(mutation)
          //
          // }
          switch (mutation.type) {
            case "childList":
              array = [];
              mutation.addedNodes.forEach((n) => {
                if (n.innerHTML !== "") {
                  array.push(n.innerHTML);
                }
              });
              if (array.length !== 0) {
                ProxyLog("muatation.childList", {
                  href: document.location.href,
                  childNodes: array,
                });
              }
              break;
            case "attributes":
              // console.log(mutation);
              // ProxyLog("muatation.attributes", {
              //   name: mutation.attributeName,
              //   oldValue: mutation.oldValue || "",
              //   newValue: mutation.target[mutation.attributeName] || "",
              //   // innerHTML: mutation.target.innerHTML,
              // });
              break;
            default:
          }
        });
      });

      mutationObserver.observe(document.documentElement, {
        href: document.location.href,
        attributes: true,
        characterData: true,
        childList: true,
        subtree: true,
        attributeOldValue: true,
        characterDataOldValue: true,
      });
    });

    function logPostMessageEvents(event) {
      try {
        ProxyLog("postMessage", {
          href: document.location.href,
          origin: event.origin,
          source: event.source.location.href,
          data: event.data,
        });
        // console.log({
        //     "byHref": document.location.href,
        //     "origin": event.origin,
        //     "source": event.source.location.href,
        //     "data":   event.data,
        // });

        // console.log("Message received by: " + document.location.href, "\norigin: " + event.origin + " source: ", event.source, "\ndata:", event.data)
      } catch (error) {
        // If the source window is cross-origin, you can't log it here
        ProxyLog("postMessage", {
          docLocationHref: document.location.href,
          origin: event.origin,
          source: event.source.location.href,
          data: event.data,
        });
      }
    }
    addEventListener("message", logPostMessageEvents);
  }
})();
