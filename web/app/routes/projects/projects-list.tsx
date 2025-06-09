import type { Project } from "~/services/user";
import ProjectCard from "./components/project-card";
import { Link } from "react-router";

type ProjectsListProps = {
  projects: Project[];
  isLoading?: boolean;
};

const ProjectsList = ({ projects, isLoading = false }: ProjectsListProps) => {
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {[...Array(4)].map((_, index) => (
          <div
            key={index}
            className="rounded-md border border-gray-200 bg-white p-6 dark:border-gray-700 dark:bg-gray-800"
          >
            <div className="animate-pulse">
              <div className="h-5 bg-gray-200 rounded w-3/4 mb-4 dark:bg-gray-700"></div>
              <div className="space-y-3">
                <div className="h-4 bg-gray-200 rounded w-1/2 dark:bg-gray-700"></div>
                <div className="h-4 bg-gray-200 rounded w-2/3 dark:bg-gray-700"></div>
              </div>
              <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
                <div className="h-3 bg-gray-200 rounded w-full dark:bg-gray-700"></div>
              </div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (projects.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 px-4">
        <svg
          className="h-12 w-12 text-gray-400 mb-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
          />
        </svg>
        <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
          No projects yet
        </h3>
        <p className="text-sm text-gray-500 dark:text-gray-400 text-center max-w-sm">
          Get started by creating your first project. Projects help you organize
          your work and collaborate with your team.
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 p-4">
      {projects.map((project) => (
        <Link to={`/projects/${project.uuid}/dashboard`}>
          <ProjectCard key={project.uuid} project={project} />
        </Link>
      ))}
    </div>
  );
};

export default ProjectsList;
