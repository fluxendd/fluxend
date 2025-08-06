import type { Project } from "~/services/user";
import ProjectCard from "./components/project-card";
import { Link, useNavigate } from "react-router";

type ProjectsListProps = {
  projects: Project[];
  isLoading?: boolean;
};

const ProjectsList = ({ projects, isLoading = false }: ProjectsListProps) => {
  const navigate = useNavigate();
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {[...Array(4)].map((_, index) => (
          <div
            key={index}
            className="rounded-lg border border-gray-200 bg-white p-6 dark:border-gray-700 dark:bg-gray-800"
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

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 p-4">
      {projects.map((project) => (
        <Link key={project.uuid} to={`/projects/${project.uuid}/dashboard`}>
          <ProjectCard 
            project={project} 
            onDocsClick={() => navigate(`/projects/${project.uuid}/docs`)}
          />
        </Link>
      ))}
    </div>
  );
};

export default ProjectsList;
