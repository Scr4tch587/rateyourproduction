"use client";

import { useState, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/components/auth-provider";
import { apiFetch } from "@/lib/api";
import type { CreateLogRequest, Production } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface LogFormProps {
  workId: string;
  productions: Production[];
  defaultProductionId?: string;
}

export function LogForm({ workId, productions, defaultProductionId }: LogFormProps) {
  const router = useRouter();
  const { user, accessToken } = useAuth();
  const [productionId, setProductionId] = useState(defaultProductionId || "");
  const [seenDate, setSeenDate] = useState("");
  const [rating, setRating] = useState("4.0");
  const [reviewText, setReviewText] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!user) {
      setError("Log in to add a record.");
      return;
    }
    if (!accessToken) {
      setError("No Supabase session available.");
      return;
    }

    setSubmitting(true);
    setError("");
    const payload: CreateLogRequest = {
      work_id: workId,
      production_id: productionId || undefined,
      seen_date: seenDate || undefined,
      rating: rating ? Number(rating) : undefined,
      review_text: reviewText || undefined,
    };

    try {
      await apiFetch("/api/v1/logs", {
        method: "POST",
        body: JSON.stringify(payload),
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });
      setReviewText("");
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save log");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-2 rounded-sm border border-border bg-card p-3">
      <div className="font-serif text-sm font-bold">Log This Work</div>
      <div className="grid gap-2 md:grid-cols-2">
        <label className="text-[11px] text-muted-foreground">
          Production
          <select
            value={productionId}
            onChange={(e) => setProductionId(e.target.value)}
            className="mt-1 h-8 w-full border border-input bg-background px-2 text-xs"
          >
            <option value="">No specific production</option>
            {productions.map((production) => (
              <option key={production.id} value={production.id}>
                {production.production_label || production.company_name || production.city || production.slug}
              </option>
            ))}
          </select>
        </label>
        <label className="text-[11px] text-muted-foreground">
          Seen date
          <Input
            type="date"
            value={seenDate}
            onChange={(e) => setSeenDate(e.target.value)}
            className="mt-1 h-8 bg-background text-xs"
          />
        </label>
      </div>
      <label className="text-[11px] text-muted-foreground block">
        Rating
        <Input
          type="number"
          step="0.5"
          min="0.5"
          max="5"
          value={rating}
          onChange={(e) => setRating(e.target.value)}
          className="mt-1 h-8 bg-background text-xs"
        />
      </label>
      <label className="text-[11px] text-muted-foreground block">
        Review
        <textarea
          value={reviewText}
          onChange={(e) => setReviewText(e.target.value)}
          rows={4}
          className="mt-1 w-full border border-input bg-background px-2 py-2 text-xs outline-none"
          placeholder="Optional notes on the performance"
        />
      </label>
      {error && <div className="text-xs text-destructive">{error}</div>}
      <Button type="submit" size="sm" disabled={submitting}>
        {submitting ? "Saving..." : "Save Log"}
      </Button>
    </form>
  );
}
