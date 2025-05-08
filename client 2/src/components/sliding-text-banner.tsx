import Marquee from 'react-fast-marquee';
import * as Constants from "@/constants";
import { Star3 } from "@/components/custom-icons";

export default function SlidingTextBanner({ CustomIcon = Star3, text = Constants.SLIDING_TEXT }: { CustomIcon?: React.FC; text?: string; }) {
  return (
    <div>
      <Marquee
        className="border-y-border dark:border-y-darkBorder dark:border-darkBorder dark:bg-secondaryBlack border-y-2 bg-white py-3 font-base sm:py-5"
        direction="left"
      >
        {Array(10)
          .fill('xd')
          .map((x, id) => {
            return (
              <div className="flex items-center" key={id}>
                <span className="mx-10 text-xl font-heading sm:text-2xl lg:text-4xl">
                  { text }
                </span>
                <CustomIcon />
              </div>
            )
          })}
      </Marquee>
    </div>
  )
}
