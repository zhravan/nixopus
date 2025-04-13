'use client';
import { LoginForm } from '@/components/features/login-form';
import { useTranslation } from '@/hooks/use-translation';
import { useLoginUserMutation, useTwoFactorLoginMutation } from '@/redux/services/users/authApi';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { toast } from 'sonner';

export default function LoginPage() {
  const { t } = useTranslation();
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [code, setCode] = useState('');
  const [showTwoFactor, setShowTwoFactor] = useState(false);
  const [loginUser, { isLoading: isLoginLoading }] = useLoginUserMutation();
  const [twoFactorLogin, { isLoading: isTwoFactorLoading }] = useTwoFactorLoginMutation();

  const handleLogin = async () => {
    try {
      const response = await loginUser({ email, password }).unwrap();
      console.log(
        "response",
        response
      );
      if (response.user.two_factor_enabled) {
        setShowTwoFactor(true);
      } else {
        router.push('/dashboard');
      }
    } catch (error) {
      toast.error(t('auth.login.errors.loginFailed'));
    }
  };

  const handleTwoFactorLogin = async () => {
    try {
      await twoFactorLogin({ email, password, code }).unwrap();
      router.push('/dashboard');
    } catch (error) {
      toast.error(t('auth.login.errors.2faFailed'));
    }
  };

  return (
    <div className="flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm md:max-w-3xl">
        <LoginForm
          email={email}
          password={password}
          handleEmailChange={(e) => setEmail(e.target.value)}
          handlePasswordChange={(e) => setPassword(e.target.value)}
          handleLogin={handleLogin}
          isLoading={isLoginLoading}
          twoFactorCode={code}
          handleTwoFactorCodeChange={(e) => setCode(e.target.value)}
          showTwoFactor={showTwoFactor}
          handleTwoFactorLogin={handleTwoFactorLogin}
          isTwoFactorLoading={isTwoFactorLoading}
        />
      </div>
    </div>
  );
}
