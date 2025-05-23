type AppHeaderProps = {
  title: string;
};

export function AppHeader({ title }: AppHeaderProps) {
  if (!title) {
    return null;
  }

  return (
    <header className="group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 flex h-12 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear">
      <h1 className="text-base font-medium">{title}</h1>
    </header>
  );
}
