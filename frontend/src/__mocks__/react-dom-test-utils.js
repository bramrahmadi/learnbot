// Custom mock for react-dom/test-utils to fix React 19 compatibility
// React 19 removed act from react-dom/test-utils and moved it to React.act
// but React 19.2.4 doesn't export act from the main react package in jsdom
// This provides a working act implementation using react-dom/client

const ReactDOM = require("react-dom/client");

// Use the act from react-dom/client if available, otherwise provide a simple implementation
let actImpl;

try {
  // Try to get act from react-dom/client
  const { act } = require("react-dom/client");
  actImpl = act;
} catch (e) {
  // Fallback: simple synchronous act
  actImpl = function act(callback) {
    const result = callback();
    if (result && typeof result.then === "function") {
      return result;
    }
    return {
      then: function(resolve) {
        resolve(result);
        return { then: function(r) { if (r) r(undefined); return { then: function() {} }; } };
      }
    };
  };
}

module.exports = { act: actImpl };
