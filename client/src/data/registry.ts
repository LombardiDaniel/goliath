const UI = [
  {
    name: "accordion",
    title: "Accordion",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-accordion"],
    files: [
      {
        path: "src/components/ui/accordion.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "alert",
    title: "Alert",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/alert.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "alert-dialog",
    title: "Alert dialog",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-alert-dialog"],
    registryDependencies: ["https://neobrutalism.dev/r/nbutton.json"],
    files: [
      {
        path: "src/components/ui/alert-dialog.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "avatar",
    title: "Avatar",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-avatar"],
    files: [
      {
        path: "src/components/ui/avatar.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "badge",
    title: "Badge",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/badge.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "breadcrumb",
    title: "Breadcrumb",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-slot"],
    files: [
      {
        path: "src/components/ui/breadcrumb.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "nbutton",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-slot"],
    files: [
      {
        path: "src/components/ui/button.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "button",
    title: "Button",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-slot"],
    files: [
      {
        path: "src/components/ui/button.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "calendar",
    title: "Calendar",
    type: "registry:ui",
    dependencies: ["react-day-picker@8.10.1", "date-fns"],
    registryDependencies: ["https://neobrutalism.dev/r/nbutton.json"],
    files: [
      {
        path: "src/components/ui/calendar.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "ncard",
    title: "Card",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/card.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "card",
    title: "Card",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/card.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "carousel",
    title: "Carousel",
    type: "registry:ui",
    dependencies: ["embla-carousel-react"],
    registryDependencies: ["https://neobrutalism.dev/r/nbutton.json"],
    files: [
      {
        path: "src/components/ui/carousel.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "chart",
    title: "Chart",
    type: "registry:ui",
    dependencies: ["recharts", "lucide-react"],
    registryDependencies: ["https://neobrutalism.dev/r/ncard.json"],
    files: [
      {
        path: "src/components/ui/chart.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "checkbox",
    title: "Checkbox",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-checkbox"],
    files: [
      {
        path: "src/components/ui/checkbox.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "collapsible",
    title: "Collapsible",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-collapsible"],
    files: [
      {
        path: "src/components/ui/collapsible.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "command",
    title: "Command",
    type: "registry:ui",
    dependencies: ["cmdk"],
    registryDependencies: ["https://neobrutalism.dev/r/ndialog.json"],
    files: [
      {
        path: "src/components/ui/command.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "context-menu",
    title: "Context menu",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-context-menu"],
    files: [
      {
        path: "src/components/ui/context-menu.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "ndialog",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-dialog"],
    files: [
      {
        path: "src/components/ui/dialog.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "dialog",
    title: "Dialog",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-dialog"],
    files: [
      {
        path: "src/components/ui/dialog.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "drawer",
    title: "Drawer",
    type: "registry:ui",
    dependencies: ["vaul", "@radix-ui/react-dialog"],
    files: [
      {
        path: "src/components/ui/drawer.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "dropdown-menu",
    title: "Dropdown menu",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-dropdown-menu"],
    files: [
      {
        path: "src/components/ui/dropdown-menu.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "form",
    title: "Form",
    type: "registry:ui",
    dependencies: [
      "@radix-ui/react-label",
      "@radix-ui/react-slot",
      "@hookform/resolvers",
      "zod",
      "react-hook-form",
    ],
    registryDependencies: [
      "https://neobrutalism.dev/r/nbutton.json",
      "https://neobrutalism.dev/r/nlabel.json",
    ],
    files: [
      {
        path: "src/components/ui/form.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "hover-card",
    title: "Hover card",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-hover-card"],
    files: [
      {
        path: "src/components/ui/hover-card.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "image-card",
    title: "Image card",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/image-card.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "input",
    title: "Input",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/input.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "ninput",
    title: "Input",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/input.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "input-otp",
    title: "Input OTP",
    type: "registry:ui",
    dependencies: ["input-otp"],
    files: [
      {
        path: "src/components/ui/input-otp.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "nlabel",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-label"],
    files: [
      {
        path: "src/components/ui/label.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "label",
    title: "Label",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-label"],
    files: [
      {
        path: "src/components/ui/label.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "marquee",
    title: "Marquee",
    docs: "Make sure there is enough content in `items` so it loops perfectly. Visit https://jackwhiting.co.uk/posts/creating-a-marquee-with-tailwind to learn more.",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/marquee.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "menubar",
    title: "Menubar",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-menubar"],
    files: [
      {
        path: "src/components/ui/menubar.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "navigation-menu",
    title: "Navigation menu",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-navigation-menu"],
    files: [
      {
        path: "src/components/ui/navigation-menu.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "pagination",
    title: "Pagination",
    type: "registry:ui",
    registryDependencies: ["https://neobrutalism.dev/r/nbutton.json"],
    files: [
      {
        path: "src/components/ui/pagination.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "popover",
    title: "Popover",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-popover"],
    files: [
      {
        path: "src/components/ui/popover.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "progress",
    title: "Progress",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-progress"],
    files: [
      {
        path: "src/components/ui/progress.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "radio-group",
    title: "Radio group",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-radio-group"],
    files: [
      {
        path: "src/components/ui/radio-group.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "resizable",
    title: "Resizable",
    type: "registry:ui",
    dependencies: ["react-resizable-panels"],
    files: [
      {
        path: "src/components/ui/resizable.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "scroll-area",
    title: "Scroll area",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-scroll-area"],
    files: [
      {
        path: "src/components/ui/scroll-area.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "select",
    title: "Select",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-select"],
    files: [
      {
        path: "src/components/ui/select.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "sheet",
    title: "Sheet",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-dialog"],
    files: [
      {
        path: "src/components/ui/sheet.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "nsheet",
    title: "Sheet",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-dialog"],
    files: [
      {
        path: "src/components/ui/sheet.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "sidebar",
    title: "Sidebar",
    type: "registry:ui",
    registryDependencies: [
      "https://neobrutalism.dev/r/nbutton.json",
      "https://neobrutalism.dev/r/nsheet.json",
      "https://neobrutalism.dev/r/ntooltip.json",
      "https://neobrutalism.dev/r/ninput.json",
      "https://neobrutalism.dev/r/nskeleton.json",
    ],
    files: [
      {
        path: "src/components/ui/sidebar.tsx",
        type: "registry:ui",
      },
      {
        path: "src/hooks/use-mobile.ts",
        type: "registry:hook",
      },
    ],
  },
  {
    name: "skeleton",
    title: "Skeleton",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/skeleton.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "nskeleton",
    title: "Skeleton",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/skeleton.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "slider",
    title: "Slider",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-slider"],
    files: [
      {
        path: "src/components/ui/slider.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "sonner",
    title: "Sonner",
    type: "registry:ui",
    dependencies: ["sonner", "next-themes"],
    files: [
      {
        path: "src/components/ui/sonner.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "switch",
    title: "Switch",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-switch"],
    files: [
      {
        path: "src/components/ui/switch.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "table",
    title: "Table",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/table.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "tabs",
    title: "Tabs",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-tabs"],
    files: [
      {
        path: "src/components/ui/tabs.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "textarea",
    title: "Textarea",
    type: "registry:ui",
    files: [
      {
        path: "src/components/ui/textarea.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "tooltip",
    title: "Tooltip",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-tooltip"],
    files: [
      {
        path: "src/components/ui/tooltip.tsx",
        type: "registry:ui",
      },
    ],
  },
  {
    name: "ntooltip",
    title: "Tooltip",
    type: "registry:ui",
    dependencies: ["@radix-ui/react-tooltip"],
    files: [
      {
        path: "src/components/ui/tooltip.tsx",
        type: "registry:ui",
      },
    ],
  },
]

const STARS = [
  {
    name: "s1",
    title: "Star 1",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s1.tsx",
        type: "registry:component",
        target: "components/stars/s1.tsx",
      },
    ],
  },
  {
    name: "s2",
    title: "Star 2",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s2.tsx",
        type: "registry:component",
        target: "components/stars/s2.tsx",
      },
    ],
  },
  {
    name: "s3",
    title: "Star 3",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s3.tsx",
        type: "registry:component",
        target: "components/stars/s3.tsx",
      },
    ],
  },
  {
    name: "s4",
    title: "Star 4",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s4.tsx",
        type: "registry:component",
        target: "components/stars/s4.tsx",
      },
    ],
  },
  {
    name: "s5",
    title: "Star 5",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s5.tsx",
        type: "registry:component",
        target: "components/stars/s5.tsx",
      },
    ],
  },
  {
    name: "s6",
    title: "Star 6",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s6.tsx",
        type: "registry:component",
        target: "components/stars/s6.tsx",
      },
    ],
  },
  {
    name: "s7",
    title: "Star 7",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s7.tsx",
        type: "registry:component",
        target: "components/stars/s7.tsx",
      },
    ],
  },
  {
    name: "s8",
    title: "Star 8",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s8.tsx",
        type: "registry:component",
        target: "components/stars/s8.tsx",
      },
    ],
  },
  {
    name: "s9",
    title: "Star 9",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s9.tsx",
        type: "registry:component",
        target: "components/stars/s9.tsx",
      },
    ],
  },
  {
    name: "s10",
    title: "Star 10",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s10.tsx",
        type: "registry:component",
        target: "components/stars/s10.tsx",
      },
    ],
  },
  {
    name: "s11",
    title: "Star 11",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s11.tsx",
        type: "registry:component",
        target: "components/stars/s11.tsx",
      },
    ],
  },
  {
    name: "s12",
    title: "Star 12",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s12.tsx",
        type: "registry:component",
        target: "components/stars/s12.tsx",
      },
    ],
  },
  {
    name: "s13",
    title: "Star 13",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s13.tsx",
        type: "registry:component",
        target: "components/stars/s13.tsx",
      },
    ],
  },
  {
    name: "s14",
    title: "Star 14",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s14.tsx",
        type: "registry:component",
        target: "components/stars/s14.tsx",
      },
    ],
  },
  {
    name: "s15",
    title: "Star 15",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s15.tsx",
        type: "registry:component",
        target: "components/stars/s15.tsx",
      },
    ],
  },
  {
    name: "s16",
    title: "Star 16",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s16.tsx",
        type: "registry:component",
        target: "components/stars/s16.tsx",
      },
    ],
  },
  {
    name: "s17",
    title: "Star 17",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s17.tsx",
        type: "registry:component",
        target: "components/stars/s17.tsx",
      },
    ],
  },
  {
    name: "s18",
    title: "Star 18",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s18.tsx",
        type: "registry:component",
        target: "components/stars/s18.tsx",
      },
    ],
  },
  {
    name: "s19",
    title: "Star 19",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s19.tsx",
        type: "registry:component",
        target: "components/stars/s19.tsx",
      },
    ],
  },
  {
    name: "s20",
    title: "Star 20",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s20.tsx",
        type: "registry:component",
        target: "components/stars/s20.tsx",
      },
    ],
  },
  {
    name: "s21",
    title: "Star 21",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s21.tsx",
        type: "registry:component",
        target: "components/stars/s21.tsx",
      },
    ],
  },
  {
    name: "s22",
    title: "Star 22",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s22.tsx",
        type: "registry:component",
        target: "components/stars/s22.tsx",
      },
    ],
  },
  {
    name: "s23",
    title: "Star 23",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s23.tsx",
        type: "registry:component",
        target: "components/stars/s23.tsx",
      },
    ],
  },
  {
    name: "s24",
    title: "Star 24",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s24.tsx",
        type: "registry:component",
        target: "components/stars/s24.tsx",
      },
    ],
  },
  {
    name: "s25",
    title: "Star 25",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s25.tsx",
        type: "registry:component",
        target: "components/stars/s25.tsx",
      },
    ],
  },
  {
    name: "s26",
    title: "Star 26",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s26.tsx",
        type: "registry:component",
        target: "components/stars/s26.tsx",
      },
    ],
  },
  {
    name: "s27",
    title: "Star 27",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s27.tsx",
        type: "registry:component",
        target: "components/stars/s27.tsx",
      },
    ],
  },
  {
    name: "s28",
    title: "Star 28",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s28.tsx",
        type: "registry:component",
        target: "components/stars/s28.tsx",
      },
    ],
  },
  {
    name: "s29",
    title: "Star 29",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s29.tsx",
        type: "registry:component",
        target: "components/stars/s29.tsx",
      },
    ],
  },
  {
    name: "s30",
    title: "Star 30",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s30.tsx",
        type: "registry:component",
        target: "components/stars/s30.tsx",
      },
    ],
  },
  {
    name: "s31",
    title: "Star 31",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s31.tsx",
        type: "registry:component",
        target: "components/stars/s31.tsx",
      },
    ],
  },
  {
    name: "s32",
    title: "Star 32",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s32.tsx",
        type: "registry:component",
        target: "components/stars/s32.tsx",
      },
    ],
  },
  {
    name: "s33",
    title: "Star 33",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s33.tsx",
        type: "registry:component",
        target: "components/stars/s33.tsx",
      },
    ],
  },
  {
    name: "s34",
    title: "Star 34",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s34.tsx",
        type: "registry:component",
        target: "components/stars/s34.tsx",
      },
    ],
  },
  {
    name: "s35",
    title: "Star 35",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s35.tsx",
        type: "registry:component",
        target: "components/stars/s35.tsx",
      },
    ],
  },
  {
    name: "s36",
    title: "Star 36",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s36.tsx",
        type: "registry:component",
        target: "components/stars/s36.tsx",
      },
    ],
  },
  {
    name: "s37",
    title: "Star 37",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s37.tsx",
        type: "registry:component",
        target: "components/stars/s37.tsx",
      },
    ],
  },
  {
    name: "s38",
    title: "Star 38",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s38.tsx",
        type: "registry:component",
        target: "components/stars/s38.tsx",
      },
    ],
  },
  {
    name: "s39",
    title: "Star 39",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s39.tsx",
        type: "registry:component",
        target: "components/stars/s39.tsx",
      },
    ],
  },
  {
    name: "s40",
    title: "Star 40",
    type: "registry:component",
    files: [
      {
        path: "src/components/stars/s40.tsx",
        type: "registry:component",
        target: "components/stars/s40.tsx",
      },
    ],
  },
]

const REGISTRY = [...UI, ...STARS]

export default REGISTRY
