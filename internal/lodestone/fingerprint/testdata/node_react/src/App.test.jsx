import { describe, it, expect } from "vitest";
import App from "./App.jsx";

describe("App", () => {
  it("renders", () => {
    expect(typeof App).toBe("function");
  });
});
