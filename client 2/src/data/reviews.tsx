const reviewersArr = [
  {
    fullName: 'Werner Heisenberg',
    jobTitle: 'Uncertainty Principle Specialist',
    pfp: 'https://avatars.githubusercontent.com/u/17596936',
    review: 'This platform is complete garbage.',
  },
  {
    fullName: 'Niels Bohr',
    jobTitle: 'Quantum Mechanics Guru',
    pfp: 'https://avatars.githubusercontent.com/u/38118870',
    review: "I don't believe there are people who will actually use this lmao.",
  },
  {
    fullName: 'Erwin Schrödinger',
    jobTitle: 'Wave Function Wizard',
    pfp: 'https://avatars.githubusercontent.com/u/241352',
    review: 'This platform frankly sucks. Unbelievable somebody actually spent time on this.',
  },
  {
    fullName: 'Karl Schwarzschild',
    jobTitle: 'Black Hole Theorist',
    pfp: 'https://avatars.githubusercontent.com/u/16783648',
    review: 'If you want to bankrupt your company, then this platform is perfect for you.',
  },
  {
    fullName: 'Louis de Broglie',
    jobTitle: 'Wave-Particle Duality Enthusiast',
    pfp: 'https://avatars.githubusercontent.com/u/241376',
    review: 'I thought FTX had the worst risk management.',
  },
  {
    fullName: 'Albert Einstein',
    jobTitle: 'Relativity Rockstar',
    pfp: 'https://avatars.githubusercontent.com/u/241398',
    review: 'I want to vomit.',
  },
  {
    fullName: 'Paul Dirac',
    jobTitle: 'Quantum Field Theorist',
    pfp: 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/696.jpg',
    review: 'HAHAHAHAHAHAHAHAHAHAHAHAHAHAH IMAGINE USING THIS.',
  },
]

// export default reviewers
export function getReviews(count: number) {
  // return reviewersArr
  // Check if array is empty or count is invalid
  if (!reviewersArr.length || count <= 0 || count > reviewersArr.length) {
    return [];
  }

  const shuffled = [...reviewersArr];

  // Fisher-Yates shuffle algorithm
  for (let i = shuffled.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
  }

  // Return the first N elements after shuffling
  return shuffled.slice(0, count);
}