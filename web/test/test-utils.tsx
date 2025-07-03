import { ReactElement, ReactNode } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ThemeProvider } from '~/contexts/theme-context';

// Create a custom render function that includes providers
interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  initialEntries?: string[];
  queryClient?: QueryClient;
}

// Create a test query client with default options suitable for tests
function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        gcTime: 0,
        staleTime: 0,
      },
      mutations: {
        retry: false,
      },
    },
  });
}

function AllTheProviders({ 
  children,
  queryClient = createTestQueryClient(),
}: { 
  children: ReactNode;
  queryClient?: QueryClient;
}) {
  return (
    <MemoryRouter>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider>
          {children}
        </ThemeProvider>
      </QueryClientProvider>
    </MemoryRouter>
  );
}

function customRender(
  ui: ReactElement,
  {
    initialEntries = ['/'],
    queryClient,
    ...renderOptions
  }: CustomRenderOptions = {}
) {
  const testQueryClient = queryClient || createTestQueryClient();

  const Wrapper = ({ children }: { children: ReactNode }) => (
    <MemoryRouter initialEntries={initialEntries}>
      <QueryClientProvider client={testQueryClient}>
        <ThemeProvider>
          {children}
        </ThemeProvider>
      </QueryClientProvider>
    </MemoryRouter>
  );

  return {
    ...render(ui, { wrapper: Wrapper, ...renderOptions }),
    queryClient: testQueryClient,
  };
}

// Re-export everything
export * from '@testing-library/react';
export { customRender as render, createTestQueryClient };