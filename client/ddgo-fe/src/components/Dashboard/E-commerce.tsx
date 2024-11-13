"use client";
import dynamic from "next/dynamic";
import React from "react";
import ChartOne from "../Charts/ChartOne";
import ChartTwo from "../Charts/ChartTwo";
import ChatCard from "../Chat/ChatCard";
import TableOne from "../Tables/TableOne";

const MapOne = dynamic(() => import("@/components/Maps/MapOne"), {
  ssr: false,
});

const ChartThree = dynamic(() => import("@/components/Charts/ChartThree"), {
  ssr: false,
});

const ECommerce: React.FC = () => {
  return (
    <>
      <div className="mt-4 grid grid-cols-12 gap-4 md:mt-6 md:gap-6 2xl:mt-7.5 2xl:gap-7.5">
        {/* Each chart spans the full width of the grid row */}
        <div className="col-span-12">
          <ChartOne />
        </div>
        
        <div className="col-span-12">
          <ChartTwo />
        </div>

        <div className="col-span-12">
          <ChartThree />
        </div>

        <div className="col-span-12">
          <MapOne />
        </div>

        <div className="col-span-12">
          <TableOne />
        </div>

        <div className="col-span-12">
          <ChatCard />
        </div>
      </div>
    </>
  );
};

export default ECommerce;
