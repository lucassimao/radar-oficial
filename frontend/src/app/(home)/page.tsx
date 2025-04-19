import Features from "@/components/home/Features";
import Footer from "@/components/home/Footer";
import Header from "@/components/home/Header";
import Hero from "@/components/home/Hero";
import Pricing from "@/components/home/Pricing";
import Showcase from "@/components/home/Showcase";

export default function Home() {
  return (
    <main className="min-h-screen">
      <Header />
      <Hero />
      <Features />
      <Showcase />
      <Pricing />
      <Footer />
    </main>
  );
}
