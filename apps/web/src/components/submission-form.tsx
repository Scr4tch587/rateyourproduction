"use client";

import { useState, type FormEvent } from "react";
import { useAuth } from "@/components/auth-provider";
import { apiFetch } from "@/lib/api";
import type { SubmissionRequest } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function SubmissionForm({ workId }: { workId: string }) {
  const { user, accessToken } = useAuth();
  const [productionLabel, setProductionLabel] = useState("");
  const [city, setCity] = useState("");
  const [country, setCountry] = useState("");
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [message, setMessage] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!user) {
      setMessage("Log in to submit a production.");
      return;
    }
    if (!accessToken) {
      setMessage("No Supabase session available.");
      return;
    }
    setSubmitting(true);
    setMessage("");
    const payload: SubmissionRequest = {
      work_id: workId,
      production_label: productionLabel || undefined,
      city: city || undefined,
      country: country || undefined,
      start_date: startDate || undefined,
      end_date: endDate || undefined,
    };
    try {
      await apiFetch("/api/v1/submissions", {
        method: "POST",
        body: JSON.stringify(payload),
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });
      setProductionLabel("");
      setCity("");
      setCountry("");
      setStartDate("");
      setEndDate("");
      setMessage("Submitted for review.");
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "Submission failed");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-2 rounded-sm border border-border bg-card p-3">
      <div className="font-serif text-sm font-bold">Add Missing Production</div>
      <Input
        value={productionLabel}
        onChange={(e) => setProductionLabel(e.target.value)}
        placeholder="Production label"
        className="bg-background text-xs"
      />
      <div className="grid gap-2 md:grid-cols-2">
        <Input
          value={city}
          onChange={(e) => setCity(e.target.value)}
          placeholder="City"
          className="bg-background text-xs"
        />
        <Input
          value={country}
          onChange={(e) => setCountry(e.target.value)}
          placeholder="Country"
          className="bg-background text-xs"
        />
      </div>
      <div className="grid gap-2 md:grid-cols-2">
        <Input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} className="bg-background text-xs" />
        <Input type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} className="bg-background text-xs" />
      </div>
      {message && <div className="text-xs text-muted-foreground">{message}</div>}
      <Button type="submit" size="sm" disabled={submitting}>
        {submitting ? "Submitting..." : "Submit Production"}
      </Button>
    </form>
  );
}
