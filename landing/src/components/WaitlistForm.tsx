"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { cn } from "./lib/utils";

export const WaitlistForm = () => {
  const [email, setEmail] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);

    try {
      const response = await fetch("/api/waitlist", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email }),
      });

      const data = await response.json();

      if (response.ok) {
        setIsSuccess(true);
        setEmail("");
      } else {
        setError(data.error || "Failed to join waitlist. Please try again.");
      }
    } catch (error) {
      setError("Something went wrong. Please try again later.");
      console.error("Error submitting form:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      viewport={{ once: true }}
      className="relative z-10 px-4 sm:px-6"
    >
      <div className="max-w-2xl mx-auto bg-white/5 backdrop-blur-sm rounded-2xl p-8 sm:p-12 border border-white/10">
        <div className="text-center mb-10">
          <h2 className="text-2xl sm:text-3xl md:text-4xl font-bold tracking-tight text-white mb-6">
            Join the <span className="bg-gradient-to-r from-indigo-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">Waitlist</span>
          </h2>
          <p className="text-base sm:text-lg text-white/70">
            Be the first to know when we launch and get early access to exclusive features.
          </p>
        </div>

        {isSuccess ? (
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="text-center p-8 rounded-xl bg-emerald-500/10 border border-emerald-500/20"
          >
            <h3 className="text-lg sm:text-xl font-semibold text-emerald-400 mb-4">Thank you!</h3>
            <p className="text-sm sm:text-base text-white/70">We'll keep you updated on our launch progress.</p>
          </motion.div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Enter your email"
                required
                className={cn(
                  "w-full px-4 sm:px-6 py-3 sm:py-4 rounded-lg bg-white/5 border transition-all duration-300",
                  error ? "border-red-500/50" : "border-white/10",
                  "text-sm sm:text-base text-white placeholder:text-white/50",
                  "focus:outline-none focus:ring-2 focus:ring-indigo-500/50",
                )}
              />
              {error && (
                <motion.p
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="mt-2 text-sm text-red-400"
                >
                  {error}
                </motion.p>
              )}
            </div>
            <button
              type="submit"
              disabled={isSubmitting}
              className={cn(
                "w-full px-6 sm:px-8 py-3 sm:py-4 rounded-lg",
                "bg-gradient-to-r from-indigo-500 to-purple-500",
                "text-sm sm:text-base text-white font-medium",
                "hover:from-indigo-600 hover:to-purple-600",
                "focus:outline-none focus:ring-2 focus:ring-indigo-500/50",
                "transition-all duration-300",
                "disabled:opacity-50 disabled:cursor-not-allowed"
              )}
            >
              {isSubmitting ? "Joining..." : "Join Waitlist"}
            </button>
          </form>
        )}
      </div>
    </motion.div>
  );
}; 