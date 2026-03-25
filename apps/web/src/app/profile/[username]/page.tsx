interface ProfilePageProps {
  params: Promise<{ username: string }>;
}

export default async function ProfilePage({ params }: ProfilePageProps) {
  const { username } = await params;

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="mb-6">
        <h1 className="text-lg font-bold">{username}</h1>
        <div className="text-xs text-muted-foreground">0 logs · 0 reviews</div>
      </div>
      <section>
        <h2 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-3">
          Log History
        </h2>
        <div className="text-xs text-muted-foreground">No logs yet.</div>
      </section>
    </div>
  );
}
