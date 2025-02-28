"use client"
import { Button } from "@/components/ui/button";
import { logout } from "@/redux/features/users/authSlice";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import { useRouter } from "next/navigation";

export default function Home() {
  const router = useRouter()
  const authenticated = useAppSelector(state => state.auth.isAuthenticated)
  const dispatch = useAppDispatch()

  return (
    <div className="flex h-screen flex-col items-center justify-center">
      <div className="text-5xl">Hello, Welcome to Nixopus</div>
      {authenticated ? (
        <div className="mt-10 flex justify-center gap-20 text-2xl">
          <Button onClick={() => dispatch(logout())}>Signout</Button>
          <Button onClick={() => { }}>Home</Button>
        </div>
      ) : (
        <div className="mt-10 flex justify-center gap-20 text-2xl">
          <Button onClick={() => router.push("/login")}>Signin</Button>
        </div>
      )}
    </div>
  );
} 