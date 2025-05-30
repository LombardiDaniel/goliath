/**
 * @preview ![img](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAAAAXNSR0IArs4c6QAAA9RJREFUSEt9VU1TE0sUPRPJZGMWZiZVbCSxdIUfxB/ARlxqvigUih+AFpr3QDeAGzCBlYDAQn6AoJAEEtcUpWy1gBJ9uydJFnxjlVilZMyM1d2TSXdmiiwyM92377n3nHtvS4ZhGDAASPQP1Z/4VVs1IFmmzjamrbkp6YZhMLcOACYwdWpBO7/YYqwC0Ax4AB7HILAGXC4XmBn5kgCDrJ2DrlcgSfXQpgM7AB+ZiCJJLiwvLyMSCVMaFxYW0NnZiYqu1zNrS0/iMyDR/Tw5QalURHPzVZM1A6rqp1ns7+9TB4qioKGhAXt7e5bD9Y0NXLlyGd7zXo5sAwyAI3BzcxOhUAhTU1N4/Ogx1b2jowPpTBrl0zKlyuPxoKurC3NzcxRgfPwFnj55iq2tL2i+2izUjJBBtYomJibwpL8f0WgM2aUs8rkcItEohoeHUalUMDIygnfv8rhz5y4ikQjy+TxmZmbQ29trKxYuA5G+aDSKXD6HQFMAG+vruODzwe12UyPtj4bjo++4GWpBoVhEPB5HOpOxKo1X0MqAVokk1vilYBCFQoE6Jk5BC5oVnSx7UC6fIhgM4tu3/8UeMszClqoakG3JRRnyKyrabrdRjltbW9HY2AitXDbdilnKsoydnR18+PAer1/PYXV1FUdHR6zodRIw6V+zD8iDpPnwwQPLyN5SpOFYTTC9zBYyDVVVxezsLGKxmNUfjhpo2h+MjY0iNZqilVPzxztl7+4GNwYGB/Hs2ZClkaBBbVTU4tU0DalkEqOjY4x7K2ATgMMhNA0NDWFgYMAC4AeDQNHS0hJ6enpweHjoxE6NEVsibEHxKXg1+wrt7e01iqoZEJEJvz6fglttt9Dd3V0TWTOzEGAlyLIbu7u7WFtbw/z8PFZWVmhwtHcFkc2hxg8uwmOgqQmlUgkeWcZpmWlRDd7tlqFpZQQCAWxvb4sZmyVP7W3T1CSQdDFptGAgiI+fPkFVFBC+SXREo+PjY9wMhVAoFREJR+gwrJv4LCAngMnJSfT39SEaiyGbzYJoQ7o1mUzh9PcvPE8lkVvOIRwO03WyPzE+jn/7+qyRXr2/OACW5eetz7hxvQUvX04ikUjQA/fu3cdielEYdmTt7ds39Mz09DT+SSRAJmpLyw2hq4U+ILyf/PiBQqGIa9evWbPF7/dTRwcHB/SpKiokl8S+zUn85et/aLp4EV6vl0V69oUjakbugvRiGvH2OD2XzWToCNd1vU5cIXimAS1T8ZbjA6CiEgDirFplTms8kuM0tQz42/sMYI4Fe1NyzmxXJsvGwbOQ1xk+67b+AjeLDdSyno8tAAAAAElFTkSuQmCC)
 */
export default function Star29({
  color,
  size,
  stroke,
  strokeWidth,
  pathClassName,
  width,
  height,
  ...props
}: React.SVGProps<SVGSVGElement> & {
  color?: string
  size?: number
  stroke?: string
  pathClassName?: string
  strokeWidth?: number
}) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 200 200"
      width={size ?? width}
      height={size ?? height}
      {...props}
    >
      <path
        fill={color ?? "currentColor"}
        stroke={stroke}
        strokeWidth={strokeWidth}
        className={pathClassName}
        d="M80.895 8.558 100 60.811l18.915-52.253c2.186-6.175 11.406-3.705 10.266 2.756l-9.6 54.722 42.583-35.817c5.038-4.18 11.691 2.47 7.509 7.506l-35.74 42.657 54.75-9.69c6.464-1.14 8.935 8.075 2.757 10.26l-52.279 19.095 52.279 18.906c6.178 2.185 3.707 11.401-2.757 10.261l-54.75-9.596 35.835 42.562c4.182 5.035-2.471 11.686-7.509 7.506l-42.678-35.722 9.695 54.722c1.141 6.461-8.079 8.931-10.266 2.756l-19.105-52.253-18.915 52.253c-2.186 6.175-11.406 3.705-10.266-2.756l9.6-54.722-42.583 35.817c-5.038 4.18-11.691-2.471-7.509-7.506l35.74-42.657-54.655 9.691c-6.464 1.14-8.935-8.076-2.757-10.261l52.28-19.095L8.56 81.047c-6.178-2.185-3.707-11.4 2.757-10.26l54.75 9.595L30.232 37.82c-4.182-5.035 2.471-11.686 7.51-7.506l42.677 35.722-9.79-54.722c-1.14-6.46 7.984-8.93 10.266-2.756"
      />
    </svg>
  )
}
