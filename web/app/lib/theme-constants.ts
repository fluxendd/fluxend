export const THEME_STORAGE_KEY = "fluxend-theme";

export const THEME_VALUES = {
  LIGHT: "light",
  DARK: "dark",
  SYSTEM: "system",
} as const;

export type Theme = (typeof THEME_VALUES)[keyof typeof THEME_VALUES];