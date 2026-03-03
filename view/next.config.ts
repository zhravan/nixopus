import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  output: 'standalone',
  basePath: process.env.BASE_PATH || '',
  assetPrefix: process.env.ASSET_PREFIX || undefined,
  experimental: {
    serverActions: {
      allowedOrigins: ['dashboard.nixopus.com']
    }
  }
};

export default nextConfig;
