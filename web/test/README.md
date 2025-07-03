# Testing Guide

This project uses Vitest as the testing framework along with React Testing Library for component testing.

## Running Tests

```bash
# Run all tests once
yarn test:run

# Run tests in watch mode
yarn test:watch

# Run tests with coverage
yarn test:coverage

# Run tests with UI
yarn test:ui
```

## Writing Tests

### File Naming Convention

- Unit tests: `*.test.ts` or `*.test.tsx`
- Integration tests: `*.spec.ts` or `*.spec.tsx`
- Test files should be colocated with the source files they test

### Test Structure

```typescript
import { describe, it, expect, vi } from 'vitest';

describe('Component/Function Name', () => {
  it('should do something specific', () => {
    // Arrange
    // Act
    // Assert
  });
});
```

### Testing React Components

Use the custom render function from `test-utils.tsx`:

```typescript
import { render, screen, fireEvent } from '~/test/test-utils';
import { MyComponent } from './MyComponent';

describe('MyComponent', () => {
  it('renders correctly', () => {
    render(<MyComponent />);
    expect(screen.getByText('Expected Text')).toBeInTheDocument();
  });
});
```

### Testing Hooks

```typescript
import { renderHook, act } from '@testing-library/react';
import { useMyHook } from './useMyHook';

describe('useMyHook', () => {
  it('returns expected value', () => {
    const { result } = renderHook(() => useMyHook());
    
    act(() => {
      result.current.doSomething();
    });
    
    expect(result.current.value).toBe('expected');
  });
});
```

### Testing with React Query

The test utilities provide a custom query client for testing:

```typescript
import { render, screen, waitFor } from '~/test/test-utils';
import { MyComponent } from './MyComponent';

describe('MyComponent with React Query', () => {
  it('loads and displays data', async () => {
    render(<MyComponent />);
    
    // Wait for the query to complete
    await waitFor(() => {
      expect(screen.getByText('Loaded Data')).toBeInTheDocument();
    });
  });
});
```

### Mocking

#### Mocking modules

```typescript
vi.mock('~/services/api', () => ({
  fetchData: vi.fn().mockResolvedValue({ data: 'mocked' }),
}));
```

#### Mocking fetch

```typescript
global.fetch = vi.fn().mockResolvedValue({
  ok: true,
  json: async () => ({ data: 'mocked' }),
});
```

### Testing Best Practices

1. **Test behavior, not implementation**: Focus on what the component does, not how it does it
2. **Use semantic queries**: Prefer `getByRole`, `getByLabelText`, `getByText` over test IDs
3. **Avoid testing implementation details**: Don't test state, test the output
4. **Keep tests isolated**: Each test should be independent
5. **Use descriptive test names**: Test names should clearly describe what is being tested
6. **Follow AAA pattern**: Arrange, Act, Assert

### Common Testing Patterns

#### Testing form submission

```typescript
it('submits form with correct data', async () => {
  const onSubmit = vi.fn();
  render(<Form onSubmit={onSubmit} />);
  
  await userEvent.type(screen.getByLabelText('Name'), 'John Doe');
  await userEvent.click(screen.getByRole('button', { name: 'Submit' }));
  
  expect(onSubmit).toHaveBeenCalledWith({ name: 'John Doe' });
});
```

#### Testing async operations

```typescript
it('displays loading state then data', async () => {
  render(<AsyncComponent />);
  
  expect(screen.getByText('Loading...')).toBeInTheDocument();
  
  await waitFor(() => {
    expect(screen.getByText('Data loaded')).toBeInTheDocument();
  });
});
```

#### Testing error states

```typescript
it('displays error message on failure', async () => {
  // Mock the API to fail
  vi.mocked(api.fetchData).mockRejectedValue(new Error('Failed'));
  
  render(<Component />);
  
  await waitFor(() => {
    expect(screen.getByText('Error: Failed')).toBeInTheDocument();
  });
});
```

## Coverage

Run `yarn test:coverage` to generate a coverage report. The coverage report will be available in the `coverage` directory.

Coverage thresholds can be configured in `vitest.config.ts`.