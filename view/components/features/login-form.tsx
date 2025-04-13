import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import nixopusLogo from '@/public/nixopus_logo_transparent.png';
import { useTranslation } from '@/hooks/use-translation';

export interface LoginFormProps {
  email: string;
  password: string;
  handleEmailChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  handlePasswordChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  handleLogin: () => void;
  isLoading: boolean;
}

export function LoginForm({ ...props }: LoginFormProps) {
  const { t } = useTranslation();

  return (
    <div className={cn('flex flex-col gap-6')}>
      <Card className="overflow-hidden p-0">
        <CardContent className="grid p-0 md:grid-cols-2">
          <div className="p-6 md:p-8">
            <div className="flex flex-col gap-6">
              <div className="flex flex-col items-center text-center">
                <h1 className="text-2xl font-bold">{t('auth.login.title')}</h1>
                <p className="text-muted-foreground text-balance">{t('auth.login.description')}</p>
              </div>
              <div className="grid gap-3">
                <Label htmlFor="email">{t('auth.email')}</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder={t('auth.login.emailPlaceholder')}
                  required
                  value={props.email}
                  onChange={props.handleEmailChange}
                />
              </div>
              <div className="grid gap-3">
                {/* <div className="flex items-center">
                  <Label htmlFor="password">Password</Label>
                  <a href="#" className="ml-auto text-sm underline-offset-2 hover:underline">
                    Forgot your password?
                  </a>
                </div> */}
                <Label htmlFor="password">{t('auth.password')}</Label>
                <Input
                  id="password"
                  type="password"
                  required
                  value={props.password}
                  onChange={props.handlePasswordChange}
                />
              </div>
              <Button
                type="submit"
                className="w-full"
                onClick={props.handleLogin}
                disabled={props.isLoading}
              >
                {props.isLoading ? t('auth.login.loading') : t('auth.login.submit')}
              </Button>
              {/* <div className="text-center text-sm">
                Don&apos;t have an account?{' '}
                <a href="#" className="underline underline-offset-4">
                  Sign up
                </a>
              </div> */}
            </div>
          </div>
          <div className="bg-muted relative hidden md:block">
            <img
              src={nixopusLogo.src}
              alt="Nixopus Logo"
              className="absolute inset-0 h-full w-full object-cover"
            />
          </div>
        </CardContent>
      </Card>
      <div className="text-muted-foreground text-center text-xs text-balance">
        {t('auth.login.terms')}{' '}
        <a
          href="https://docs.nixopus.com/license"
          className="underline underline-offset-4 hover:text-primary"
        >
          {t('auth.login.termsOfService')}
        </a>{' '}
        {t('auth.login.and')}{' '}
        <a
          href="https://docs.nixopus.com/privacy-policy"
          className="underline underline-offset-4 hover:text-primary"
        >
          {t('auth.login.privacyPolicy')}
        </a>
        .
      </div>
    </div>
  );
}
