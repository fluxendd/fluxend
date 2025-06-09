import { BadgeAlert } from "lucide-react";

export const InfoMessage = ({ message }: { message: string }) => (
  <div className="flex items-center p-4 text-muted-foreground">
    <BadgeAlert className="mr-2" />
    <div className="text-md">{message}</div>
  </div>
);
