// This file runs before the test framework is set up
// It patches React to add the act function for React 19 compatibility

const React = require("react");

if (typeof React.act !== "function") {
  // Use a proper act implementation from react-dom
  // React 19 moved act to react-dom/client but we need it on React itself
  // for backward compatibility with @testing-library/react
  
  // Simple synchronous act that works for most test cases
  let actQueue = [];
  let isActing = false;
  
  React.act = function act(callback) {
    if (isActing) {
      return callback();
    }
    
    isActing = true;
    try {
      const result = callback();
      
      if (result && typeof result.then === "function") {
        // Async act
        return {
          then: function(resolve, reject) {
            result.then(
              function(value) {
                isActing = false;
                resolve(value);
              },
              function(error) {
                isActing = false;
                if (reject) reject(error);
              }
            );
          }
        };
      }
      
      isActing = false;
      return {
        then: function(resolve) {
          resolve(result);
          return { then: function(r) { if (r) r(undefined); return { then: function() {} }; } };
        }
      };
    } catch (e) {
      isActing = false;
      throw e;
    }
  };
}
