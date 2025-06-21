import { Calendar, ChevronRight, Link, View } from "lucide-react";
import type { Project } from "~/services/user";

type ProjectCardProps = {
  project: Project;
  onClick?: () => void;
};

const ProjectCard = ({ project, onClick }: ProjectCardProps) => {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const getStatusColor = (status: Project["status"]) => {
    return status === "active"
      ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300"
      : "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300";
  };

  return (
    <div
      className="group relative overflow-hidden rounded-lg border border-gray-200 bg-white p-3 shadow-sm transition-all duration-200 hover:shadow-md hover:border-muted-foreground cursor-pointer dark:border-muted dark:bg-background dark:hover:border-gray-600"
      onClick={onClick}
    >
      <div className="flex flex-col mb-3 gap-1">
        <div className="flex items-center justify-between">
          <h3 className="text-md font-semibold text-gray-900 dark:text-gray-100 group-hover:text-blue-600 transition-colors capitalize">
            {project.name}
          </h3>
          <ChevronRight className="h-4 w-4" />
        </div>
        <span
          className={`w-fit rounded-full px-1.5 text-xs ${getStatusColor(
            project.status
          )}`}
        >
          {project.status}
        </span>
      </div>
      <div className="flex items-center gap-2">
        <Calendar className="h-3 w-3" />
        <span className="text-xs text-gray-700 dark:text-gray-300">
          {formatDate(project.createdAt)}
        </span>
      </div>

      {/* Card Footer - Actions */}
      <div className="mt-3 pt-3 border-t border-muted flex items-center justify-between">
        <p className="text-xs text-gray-400 dark:text-gray-500 font-mono truncate">
          {project.uuid}
        </p>
        {/* <div className="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <button
            className="text-xs text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
            onClick={(e) => {
              e.stopPropagation();
              // Handle view action
            }}
          >
            <View />
          </button>
          <button
            className="text-xs text-gray-600 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
            onClick={(e) => {
              e.stopPropagation();
              // Handle edit action
            }}
          >
            Edit
          </button>
        </div> */}
      </div>

      <div className="absolute inset-0 bg-gradient-to-r from-transparent via-gray-100/10 to-transparent -translate-x-full group-hover:translate-x-full transition-transform duration-700 dark:via-white/5" />
    </div>
  );
};

export default ProjectCard;
