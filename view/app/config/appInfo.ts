import { getBaseUrl } from '@/redux/conf';

export const getAppInfo = async () => {
  const baseUrl = await getBaseUrl();
  const apiDomain = baseUrl.replace('://', 'TEMP').replace('/api', '').replace('TEMP', '://');

  const websiteDomain = 
    process.env.VIEW_DOMAIN || 
    process.env.NEXT_PUBLIC_WEBSITE_DOMAIN || 
    'http://localhost:3000';

  return {
    // learn more about this on https://supertokens.com/docs/thirdpartyemailpassword/appinfo
    appName: 'Nixopus',
    apiDomain,
    websiteDomain,
    apiBasePath: '/auth',
    websiteBasePath: '/auth'
  };
};
