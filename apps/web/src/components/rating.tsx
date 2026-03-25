interface RatingProps {
  rating: number;
  count: number;
}

export function Rating({ rating, count }: RatingProps) {
  const fullStars = Math.floor(rating);
  const halfStar = rating % 1 >= 0.25;
  const stars = "★".repeat(fullStars) + (halfStar ? "½" : "") + "☆".repeat(5 - fullStars - (halfStar ? 1 : 0));

  return (
    <span className="text-xs whitespace-nowrap">
      <span className="text-yellow-500">{stars}</span>{" "}
      <span className="text-muted-foreground">
        {rating.toFixed(1)} ({count})
      </span>
    </span>
  );
}
