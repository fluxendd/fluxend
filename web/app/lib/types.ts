export type APIResponse<T> = {
  ok: boolean;
  status: number;
  success: boolean;
  errors: string[] | null;
  content: T | null;
};
