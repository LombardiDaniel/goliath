import { MutableIcon } from "@/components/ui/mutableIcon";
import * as Constants from "@/constants";
import SlidingTextBanner from "@/components/sliding-text-banner";

export default function Features() {

  const features = Constants.FEATURES;

  return (
    <div>
      <section id="features" className="border-t-border dark:border-t-darkBorder dark:bg-darkBg border-t-2 bg-bg py-20 font-base lg:py-[100px]">
        <h2 className="mb-14 px-5 text-center text-2xl font-heading md:text-3xl lg:mb-20 lg:text-4xl">
          Features
        </h2>

        <div className="mx-auto grid w-container max-w-full grid-cols-1 gap-5 px-5 sm:grid-cols-2 lg:grid-cols-3">
          {features.map((feature, i) => {

            return (
              <div
                className="border-border dark:border-darkBorder dark:bg-secondaryBlack shadow-light dark:shadow-dark flex flex-col gap-3 rounded-base border-2 bg-white p-5"
                key={i}
              >
                {/* <MutableIcon name={feature.lucideIcon} color="black"/> */}

                {/* {feature.getLucideIcon} */}

                {/* <Icon /> */}

                <MutableIcon name={feature.lucideIcon} weight={3} />

                {/* <i className="bi bi-ticket-fi"/> */}

                <h4 className="mt-2 text-xl font-heading">
                  {i + 1}. {feature.title} {/* {i + 1} */}
                </h4>
                <p>{feature.text}</p>
              </div>
            )
          })}
        </div>
      </section>
      <SlidingTextBanner />
    </div>
  )
}
