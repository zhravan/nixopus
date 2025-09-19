'use client';

import { useState, useEffect, useMemo } from 'react';
import { Extension } from '../components/extension-card';

const dummyExtensions: Extension[] = [
  {
    id: '1',
    name: 'Minimal Secure Server',
    description: 'Creates a non-root user, disables root SSH login, sets up firewall rules, and installs Fail2ban.',
    author: 'Nixopus',
    icon: '🛡️',
    category: 'Security',
    rating: 4.9,
    downloads: 10320,
    isVerified: true
  },
  {
    id: '2',
    name: 'Docker Ready Server',
    description: 'Installs Docker and Docker Compose with sane defaults and firewall rules for containers.',
    author: 'Nixopus',
    icon: '🐳',
    category: 'Containers',
    rating: 4.8,
    downloads: 8920,
    isVerified: true
  },
  {
    id: '3',
    name: 'Postgres Database Setup',
    description: 'Installs and configures PostgreSQL with secure defaults and optional remote access.',
    author: 'Community',
    icon: '🐘',
    category: 'Database',
    rating: 4.6,
    downloads: 7210,
    isVerified: false
  },
  {
    id: '4',
    name: 'Nginx Reverse Proxy',
    description: 'Sets up Nginx as a reverse proxy with SSL (Let’s Encrypt) and basic rate limiting.',
    author: 'Nixopus',
    icon: '🌐',
    category: 'Web Server',
    rating: 4.7,
    downloads: 11050,
    isVerified: true
  },
  {
    id: '5',
    name: 'Auto Updates & Patching',
    description: 'Enables unattended upgrades and automatic security patches for your server.',
    author: 'Security Team',
    icon: '🔄',
    category: 'Maintenance',
    rating: 4.5,
    downloads: 5420,
    isVerified: false
  },
  {
    id: '6',
    name: 'Fail2Ban DoS Protection',
    description: 'Protects your server from brute-force and DoS attacks using Fail2ban with custom rules.',
    author: 'Community',
    icon: '🚫',
    category: 'Security',
    rating: 4.4,
    downloads: 4370,
    isVerified: false
  }
];

export function useExtensions() {
  const [extensions, setExtensions] = useState<Extension[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortConfig, setSortConfig] = useState<{ key: string; direction: 'asc' | 'desc' }>({
    key: 'popularity',
    direction: 'desc'
  });
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage] = useState(12);

  // Simulate API call
  useEffect(() => {
    const loadExtensions = async () => {
      try {
        setIsLoading(true);
        setError(null);
        
        // Simulate network delay
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        setExtensions(dummyExtensions);
      } catch (err) {
        setError('Failed to load extensions');
      } finally {
        setIsLoading(false);
      }
    };

    loadExtensions();
  }, []);

  const filteredAndSortedExtensions = useMemo(() => {
    let filtered = extensions.filter(extension =>
      extension.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      extension.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
      extension.author.toLowerCase().includes(searchTerm.toLowerCase()) ||
      extension.category.toLowerCase().includes(searchTerm.toLowerCase())
    );

    filtered.sort((a, b) => {
      const aValue = a[sortConfig.key as keyof Extension];
      const bValue = b[sortConfig.key as keyof Extension];
      
      if (typeof aValue === 'string' && typeof bValue === 'string') {
        return sortConfig.direction === 'asc' 
          ? aValue.localeCompare(bValue)
          : bValue.localeCompare(aValue);
      }
      
      if (typeof aValue === 'number' && typeof bValue === 'number') {
        return sortConfig.direction === 'asc' 
          ? aValue - bValue
          : bValue - aValue;
      }
      
      return 0;
    });

    return filtered;
  }, [extensions, searchTerm, sortConfig]);

  const paginatedExtensions = useMemo(() => {
    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    return filteredAndSortedExtensions.slice(startIndex, endIndex);
  }, [filteredAndSortedExtensions, currentPage, itemsPerPage]);

  const totalPages = Math.ceil(filteredAndSortedExtensions.length / itemsPerPage);

  const handleSearchChange = (value: string) => {
    setSearchTerm(value);
    setCurrentPage(1); // Reset to first page when searching
  };

  const handleSortChange = (key: string, direction: 'asc' | 'desc') => {
    setSortConfig({ key, direction });
    setCurrentPage(1); // Reset to first page when sorting
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const handleInstall = (extension: Extension) => {
    // TODO: Implement installation logic
    console.log('Installing extension:', extension.name);
    
    // Update the extension to show as installed
    setExtensions(prev => 
      prev.map(ext => 
        ext.id === extension.id 
          ? { ...ext, isInstalled: true }
          : ext
      )
    );
  };

  const handleViewDetails = (extension: Extension) => {
    console.log('Viewing details for extension:', extension.name);
  };

  return {
    extensions: paginatedExtensions,
    isLoading,
    error,
    searchTerm,
    sortConfig,
    currentPage,
    totalPages,
    totalExtensions: filteredAndSortedExtensions.length,
    handleSearchChange,
    handleSortChange,
    handlePageChange,
    handleInstall,
    handleViewDetails
  };
}
