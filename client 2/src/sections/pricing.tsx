import PricingPlan from '@/components/pricing-plan';
import SlidingTextBanner from '@/components/sliding-text-banner';

import * as Constants from "@/constants";

export default function Pricing() {
  return (
    <div>
      <section id="pricing" className="border-b-border dark:border-b-darkBorder dark:bg-secondaryBlack inset-0 flex w-full flex-col items-center justify-center border-b-1 bg-white bg-[linear-gradient(to_right,#80808033_1px,transparent_1px),linear-gradient(to_bottom,#80808033_1px,transparent_1px)] bg-[size:70px_70px] font-base">
        <div className="mx-auto w-container max-w-full px-5 py-20 lg:py-[100px]">
          <h2 className="mb-14 text-center text-2xl font-heading md:text-3xl lg:mb-20 lg:text-4xl">
            Planos e Preços
          </h2>
          <div className="grid grid-cols-3 gap-8 w900:mx-auto w900:w-2/3 w900:grid-cols-1 w500:w-full">
            {Constants.PRICING.map((plan, i) => {
              return (
                <PricingPlan
                  key={i}
                  planName={plan.name}
                  description={plan.desc}
                  price={plan.price}
                  perks={plan.perks}
                  mostPopular={plan.mostPopular}
                />
              )
            })}
          </div>
        </div>
      </section>
      <SlidingTextBanner />
    </div>
  )
}
