"use client"
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function Home() {
  const [authenticated, setAuthenticated] = useState(false);
  const router = useRouter()

  return (
    <div className="flex h-screen flex-col items-center justify-center">
      <div className="text-5xl">Hello, Welcome to Nixopus</div>
      {authenticated ? (
        <div className="mt-10 flex justify-center gap-20 text-2xl">
          <Button onClick={() => { }}>Signout</Button>
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