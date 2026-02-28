import "@testing-library/jest-dom";

// Fix React 19 compatibility with @testing-library/react
// React 19.2.4 doesn't export `act` from the main `react` package
// but react-dom/test-utils tries to call React.act()
// We need to patch React to add the act function from react-dom/test-utils
// by using a workaround that provides a working act implementation

// eslint-disable-next-line @typescript-eslint/no-require-imports
const React = require("react");

if (typeof React.act !== "function") {
  // Provide a simple act implementation that works for synchronous tests
  React.act = function act(callback: () => unknown) {
    const result = callback();
    return {
      then: (resolve: (value: unknown) => void) => {
        if (result && typeof (result as Promise<unknown>).then === "function") {
          return (result as Promise<unknown>).then(resolve);
        }
        resolve(result);
        return { then: (r: (value: unknown) => void) => r(undefined) };
      },
    };
  };
}
