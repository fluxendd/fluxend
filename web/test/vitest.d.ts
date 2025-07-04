/// <reference types="vitest" />
import '@testing-library/jest-dom';

declare module 'vitest' {
  interface Assertion extends jest.Matchers<void, any> {}
  interface AsymmetricMatchersContaining extends jest.AsymmetricMatchers {}
}