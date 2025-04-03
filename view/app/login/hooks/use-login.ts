import { useAppSelector } from '@/redux/hooks';
import { useLoginUserMutation } from '@/redux/services/users/authApi';
import { useRouter } from 'next/navigation';
import React, { useEffect } from 'react';
import { toast } from 'sonner';

function useLogin() {
  const [email, setEmail] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [loginUser, { isLoading, error }] = useLoginUserMutation();
  const router = useRouter();

  const authenticated = useAppSelector((state) => state.auth.isAuthenticated);
  const user = useAppSelector((state) => state.auth.user);

  useEffect(() => {
    if (authenticated && user) {
      router.push('/dashboard');
    }
    return () => {};
  }, [authenticated, user]);

  const handleEmailChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(event.target.value);
  };

  const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(event.target.value);
  };

  const handleLogin = async () => {
    try {
      await loginUser({ email, password });
      router.push('/dashboard');
    } catch (error) {
      toast.error('Login failed');
    }
  };

  return {
    email,
    password,
    handleEmailChange,
    handlePasswordChange,
    handleLogin,
    isLoading,
    error
  };
}

export default useLogin;
