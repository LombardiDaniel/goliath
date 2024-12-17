import {
    BookOpenText,
    CalendarClock,
    Glasses,
    Hourglass,
    LucideProps,
    Search,
    Ticket,
    X
} from "lucide-react"; // Update this to use Lucide React
import React, { FC } from "react";

// Icon weight mappings
const IconWeights = {
  bold: 2,
  medium: 1.5,
  regular: 1,
};

type IconWeightsAllowedValues = (typeof IconWeights)[keyof typeof IconWeights];

interface MutableIconProps extends React.HTMLAttributes<HTMLDivElement> {
  name: string | CategoryTags;
  value?: string;
  color?: string;
  size?: number;
  weight: IconWeightsAllowedValues;
  onClick?: () => void | null;
}

export const iconMap: { [key: string]: FC<LucideProps> } = {
  "ticket": Ticket,
  "search": Search,
  "book-open-text": BookOpenText,
  "calendar-clock": CalendarClock,
  "glasses": Glasses,
  "hourglass": Hourglass,
//   badgePlus: BadgePlus,
//   trash: Bomb,
//   finance: Landmark,
//   personal: User,
//   work: Briefcase,
//   learning: Library,
//   health: HeartPulse,
//   coding: SquareTerminal,
//   other: FilePlus2,
//   pencil: Pencil,
//   circleX: CircleX,
//   circle: Circle,
//   badgecheck: BadgeCheck,
  x: X,
};

export enum CategoryTags {
  WORK = "work",
  PERSONAL = "personal",
  FINANCE = "finance",
  LEARNING = "learning",
  CODING = "coding",
  HEALTH = "health",
  OTHER = "other",
}

export const MutableIcon: FC<MutableIconProps> = ({
  name,
  size = 24,
  color = "currentColor",
  weight = IconWeights.regular,
  onClick,
  style,
}) => {
  const IconComponent = iconMap[name];
  if (!IconComponent) {
    return null;
  }

  return (
    <div onClick={onClick} style={style}>
      <IconComponent color={color} size={size} strokeWidth={3} />
    </div>
  );
};