'use client';
import { LoginForm } from '@/components/features/login-form';
import useLogin from './hooks/use-login';
import { ErrorBoundary } from '@/components/ui/error-handler';
import { useTranslation } from '@/hooks/use-translation';

export default function LoginPage() {
  const { t } = useTranslation();
  const {
    email,
    password,
    handleEmailChange,
    handlePasswordChange,
    handleLogin,
    isLoading,
    error
  } = useLogin();

  return (
    <div className="flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm md:max-w-3xl">
        <LoginForm
          email={email}
          password={password}
          handleEmailChange={handleEmailChange}
          handlePasswordChange={handlePasswordChange}
          handleLogin={handleLogin}
          isLoading={isLoading}
        />
      </div>
      {error && <ErrorBoundary errors={[{ error, title: t('auth.login.errors.loginFailed') }]} />}
    </div>
  );
}
