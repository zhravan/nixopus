'use client';
import { DahboardUtilityHeader } from '@/components/dashboard-page-header';
import React from 'react';
import GitHubAppSetup from './components/github-connector/github-app-setup';
import GithubRepositories, { GithubRepositoriesSkeletonLoader } from './components/github-repositories/repository-card';
import { GithubRepository } from '@/redux/types/github';
import PaginationWrapper from '@/components/pagination';
import useGithubRepoPagination from './hooks/use_github_repo_pagination';

function page() {
  const {
    isLoading,
    connectors,
    GetGithubConnectors,
    setSelectedRepository,
    searchTerm,
    handleSearchChange,
    onSortChange,
    sortOptions,
    sortConfig,
    handlePageChange,
    currentPage,
    totalPages,
    paginatedApplications
  } = useGithubRepoPagination();

  const renderGithubRepositories = () => {
    if (isLoading) {
      return <GithubRepositoriesSkeletonLoader />;
    }

    if (paginatedApplications.length === 0 && !isLoading) {
      return <div className="text-center">No repositories found</div>;
    }
    return (
      <>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {paginatedApplications &&
            paginatedApplications?.map((repo: any) => (
              <GithubRepositories
                key={repo.id}
                {...repo}
                setSelectedRepository={setSelectedRepository}
              />
            ))}
        </div>
        <div className="mt-8 flex justify-center">
          <PaginationWrapper
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={handlePageChange}
          />
        </div>
      </>
    );
  };

  return (
    <div className="container mx-auto py-6 space-y-8 max-w-4xl">
      <DahboardUtilityHeader<GithubRepository>
        searchTerm={searchTerm}
        handleSearchChange={handleSearchChange}
        sortConfig={sortConfig}
        onSortChange={onSortChange}
        sortOptions={sortOptions}
        label="Github Repositories"
        className="mt-5 mb-5"
      />
      {connectors?.length === 0 ? (
        <GitHubAppSetup GetGithubConnectors={GetGithubConnectors} />
      ) : (
        renderGithubRepositories()
      )}
    </div>
  );
}

export default page;
