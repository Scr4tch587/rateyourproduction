"use client";

import { useCallback, useEffect, useState, type FormEvent } from "react";
import Link from "next/link";
import { useAuth } from "@/components/auth-provider";
import { apiFetch } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import type {
  PaginatedResponse,
  ProductionSubmission,
  WorkCard,
} from "@/lib/types";

export function AdminPanel() {
  const { user, loading } = useAuth();
  const [works, setWorks] = useState<WorkCard[]>([]);
  const [submissions, setSubmissions] = useState<ProductionSubmission[]>([]);
  const [title, setTitle] = useState("");
  const [type, setType] = useState("play");
  const [message, setMessage] = useState("");

  const loadData = useCallback(async () => {
    try {
      const [worksRes, submissionsRes] = await Promise.all([
        apiFetch<PaginatedResponse<WorkCard>>("/api/v1/admin/works?per_page=20"),
        apiFetch<PaginatedResponse<ProductionSubmission>>("/api/v1/admin/submissions?status=pending"),
      ]);
      setWorks(worksRes.data);
      setSubmissions(submissionsRes.data);
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "Failed to load admin data");
    }
  }, []);

  useEffect(() => {
    if (user?.is_admin) {
      const timer = window.setTimeout(() => {
        void loadData();
      }, 0);
      return () => window.clearTimeout(timer);
    }
  }, [loadData, user?.is_admin]);

  async function createWork(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setMessage("");
    try {
      await apiFetch("/api/v1/admin/works", {
        method: "POST",
        body: JSON.stringify({ title, type }),
      });
      setTitle("");
      setMessage("Work created.");
      loadData();
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "Failed to create work");
    }
  }

  async function reviewSubmission(id: string, action: "approve" | "reject") {
    await apiFetch(`/api/v1/admin/submissions/${id}/${action}`, {
      method: "POST",
      body: JSON.stringify({}),
    });
    loadData();
  }

  if (loading) {
    return <div className="max-w-6xl mx-auto px-4 py-10 text-xs">Loading...</div>;
  }

  if (!user?.is_admin) {
    return (
      <div className="max-w-6xl mx-auto px-4 py-10">
        <div className="rounded-sm border border-border bg-card p-4 text-xs text-muted-foreground">
          Admin access is required.
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="mb-6 border-b border-border pb-2">
        <h1 className="font-serif text-2xl font-bold">Admin Panel</h1>
        <p className="text-xs text-muted-foreground">
          Review submitted productions and maintain the work catalogue.
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-[320px_1fr]">
        <section className="rounded-sm border border-border bg-card p-4">
          <h2 className="font-serif text-sm font-bold mb-3">Create Work</h2>
          <form onSubmit={createWork} className="space-y-3">
            <Input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Work title"
              className="bg-background text-xs"
            />
            <select
              value={type}
              onChange={(e) => setType(e.target.value)}
              className="h-8 w-full border border-input bg-background px-2 text-xs"
            >
              <option value="play">Play</option>
              <option value="musical">Musical</option>
              <option value="opera">Opera</option>
            </select>
            <Button type="submit" size="sm">Create</Button>
          </form>
          {message && <div className="mt-3 text-xs text-muted-foreground">{message}</div>}
        </section>

        <section className="space-y-6">
          <div className="rounded-sm border border-border bg-card p-4">
            <h2 className="font-serif text-sm font-bold mb-3">Pending Submissions</h2>
            <div className="space-y-3">
              {submissions.length === 0 ? (
                <div className="text-xs text-muted-foreground italic">No pending submissions.</div>
              ) : (
                submissions.map((submission) => (
                  <div key={submission.id} className="border-b border-border pb-3 text-xs last:border-b-0">
                    <div className="font-medium">{submission.work_title}</div>
                    <div className="text-muted-foreground">
                      {submission.production_label || "Untitled production"}
                      {submission.city ? ` · ${submission.city}` : ""}
                      {submission.country ? `, ${submission.country}` : ""}
                    </div>
                    <div className="mt-2 flex gap-2">
                      <Button size="xs" onClick={() => reviewSubmission(submission.id, "approve")}>
                        Approve
                      </Button>
                      <Button size="xs" variant="outline" onClick={() => reviewSubmission(submission.id, "reject")}>
                        Reject
                      </Button>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>

          <div className="rounded-sm border border-border bg-card p-4">
            <h2 className="font-serif text-sm font-bold mb-3">Works</h2>
            <div className="space-y-2">
              {works.map((work) => (
                <div key={work.id} className="flex items-center justify-between gap-3 border-b border-border pb-2 text-xs last:border-b-0">
                  <Link href={`/works/${work.slug}`} className="font-medium hover:text-burgundy">
                    {work.title}
                  </Link>
                  <span className="text-muted-foreground">{work.type}</span>
                </div>
              ))}
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}
