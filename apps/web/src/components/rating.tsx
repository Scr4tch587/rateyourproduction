interface RatingProps {
  rating: number;
  count: number;
}

export function Rating({ rating, count }: RatingProps) {
  const fullStars = Math.floor(rating);
  const halfStar = rating % 1 >= 0.25;
  const emptyStars = 5 - fullStars - (halfStar ? 1 : 0);
  const stars = [
    ...Array.from({ length: fullStars }, () => "\u2605"),
    ...(halfStar ? ["\u2726"] : []),
    ...Array.from({ length: emptyStars }, () => "\u2606"),
  ].join(" ");

  return (
    <span className="text-xs whitespace-nowrap font-serif">
      <span className="text-burgundy tracking-[0.12em]">{stars}</span>{" "}
      <span className="font-sans text-muted-foreground">
        {rating.toFixed(1)}
        {count > 0 && <span className="text-[10px]"> ({count})</span>}
      </span>
    </span>
  );
}
