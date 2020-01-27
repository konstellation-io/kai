import React, { useRef, useEffect, useState } from 'react';
import ReactDOMServer from 'react-dom/server';
import VersionNode, { TYPES } from '../Shape/Node/Node';

import { select, event } from 'd3-selection';
import { scaleBand, ScaleBand } from 'd3-scale';
import { max, range } from 'd3-array';
import { zoom } from 'd3-zoom';
import { wrap, centerText } from '../../utils/d3';

import styles from './VersionStatusViewer.module.scss';

import { NodeStatus } from '../../graphql/models';

const MARGIN_WORKFLOW_NAMES_PERC = 0.08;
const DEFAULT_NODE_WIDTH = 120.33;
const DEFAULT_NODE_HEIGHT = 37.9;
const STROKE_WIDTH = 0.7;
const SCALE_PADDING_INNER = 0.4;
const SCALE_PADDING_OUTER = 0.1;

type Node = {
  id: string;
  name: string;
  status: NodeStatus;
  type?: string;
};

type Edge = {
  id: string;
  status?: string;
  fromNode: string;
  toNode: string;
};

export type Workflow = {
  name: string;
  nodes: Node[];
  edges: Edge[];
};

type Props = {
  width: number;
  height: number;
  margin: {
    top: number;
    right: number;
    bottom: number;
    left: number;
  };
  data: Workflow[];
  published: boolean;
};

function getNodeTextPadding(type: string) {
  switch (type) {
    case TYPES.INPUT:
      return 45;
    case TYPES.OUTPUT:
      return 15;
    default:
      return 30;
  }
}

function VersionStatusViewer({
  width,
  height,
  margin,
  data,
  published
}: Props) {
  const [shouldComponentUpdate, setShouldComponentUpdate] = useState(false);
  const container = useRef(null);
  const svg = useRef(null);

  const nodeIdToIndex: any = {};
  let g: any;
  let defs: any;
  let newWorkflows: any;
  let workflowsTag: any;
  let nodes: any;
  let newNodes: any;
  let edges: any;
  let newEdges: any;
  let edgesG: any;
  let busG: any;
  let xScale: any;
  let fontSize: any;
  let nodeHeight: any;
  let outerPadding: any;
  const ref = useRef<any>({
    workflowsG: null,
    workflows: null,
    oldWorkflows: null,
    nodesG: null,
    oldNodes: null,
    nodeWidth: null,
    nodeSizeRatio: null
  });

  const marginWorkflow = width * MARGIN_WORKFLOW_NAMES_PERC;
  const marginLeft = marginWorkflow + margin.left;
  const innerWidth = width - marginLeft - margin.right;

  const maxNodesInRow = max(data, d => d.nodes.length) || 0;
  const xDomainIndexes = range(maxNodesInRow).map(n => n.toString());

  useEffect(() => {
    cleanup();
    initialize();
    setShouldComponentUpdate(true);
  }, [width, height]);

  useEffect(() => {
    if (shouldComponentUpdate) {
      updateChart();
    }
  }, [data]);

  function cleanup() {
    select(svg.current)
      .selectAll('*')
      .remove();
  }

  function buildNodeIdToIndex() {
    data.forEach((workflow: Workflow) =>
      workflow.nodes.forEach(
        (node: Node, idx: number) => (nodeIdToIndex[node.id] = idx.toString())
      )
    );
  }

  function generateBlurFilter() {
    defs = select(svg.current).append('defs');

    const filter = defs.append('filter').attr('id', 'glow');
    filter
      .append('feGaussianBlur')
      .attr('stdDeviation', '3.5')
      .attr('result', 'coloredBlur');

    const feMerge = filter.append('feMerge');
    feMerge.append('feMergeNode').attr('in', 'coloredBlur');
    feMerge.append('feMergeNode').attr('in', 'SourceGraphic');
  }

  function updateChart() {
    setData.workflows();
    setData.nodes();
    update.workflows();
    update.nodes();
    update.workflowTags();
  }

  function initialize() {
    const svgSelection = select(svg.current);

    g = svgSelection
      .append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    generateBlurFilter();
    buildNodeIdToIndex();

    const zoomLimits: [[number, number], [number, number]] = [
      [-width / 4, -height / 4],
      [width * 1.25, height * 1.25]
    ];

    // Adds zoom
    const zoomed = () => g.attr('transform', event.transform);
    const zoomFunc = zoom()
      .scaleExtent([0.5, 7])
      .translateExtent(zoomLimits)
      .on('zoom', zoomed);
    // @ts-ignore
    svgSelection.call(zoomFunc);

    // Initialize scales
    xScale = scaleBand()
      .range([marginLeft, marginLeft + innerWidth])
      .paddingInner(SCALE_PADDING_INNER)
      .paddingOuter(SCALE_PADDING_OUTER)
      .domain(xDomainIndexes);

    ref.current.nodeWidth = xScale.bandwidth();
    ref.current.nodeSizeRatio = ref.current.nodeWidth / DEFAULT_NODE_WIDTH;
    outerPadding = xScale.step() * SCALE_PADDING_OUTER;
    nodeHeight = DEFAULT_NODE_HEIGHT * ref.current.nodeSizeRatio;
    fontSize = (ref.current.nodeWidth / 21).toFixed(1);

    // Create workflows
    ref.current.workflowsG = g.append('g').classed(styles.workflows, true);
    setData.workflows();
    create.workflows();

    // Create workflows tag
    create.workflowTags();

    // Create nodes/edges groups
    ref.current.nodesG = newWorkflows
      .append('g')
      .classed(styles.nodesG, true)
      .attr('transform', 'translate(0, 50)');
    edgesG = newWorkflows
      .append('g')
      .classed(styles.edgesG, true)
      .attr('transform', 'translate(0, 50)');

    // Create nodes
    setData.nodes();
    create.nodes();

    // Create edges
    setData.edges();
    create.edges();

    // Create Data Bus
    create.bus();
  }

  const setData = {
    workflows: function() {
      ref.current.workflows = ref.current.workflowsG
        .selectAll(`.${styles.workflow}`)
        .data(data);
    },
    nodes: function() {
      nodes = ref.current.nodesG
        .selectAll(`.${styles.node}`)
        .data((d: Workflow) => d.nodes);
    },
    edges: function() {
      edges = edgesG
        .selectAll(`.${styles.edge}`)
        .data((d: Workflow) => d.edges);
    }
  };

  const update = {
    // Pass new data to children components
    workflows: function() {
      ref.current.oldWorkflows = ref.current.workflows.select(
        `.${styles.nodesG}`
      );

      ref.current.oldNodes = ref.current.workflows
        .selectAll(`.${styles.node}`)
        .data((d: Workflow) => {
          return d.nodes;
        });
    },
    // Updates node status
    nodes: function() {
      ref.current.oldNodes.each(function(d: Node) {
        // @ts-ignore
        select(this)
          .select('g')
          .html(
            ReactDOMServer.renderToString(
              <VersionNode
                type={d.type || TYPES.DEFAULT}
                width={ref.current.nodeWidth}
                height={DEFAULT_NODE_HEIGHT * ref.current.nodeSizeRatio}
                status={d.status}
              />
            )
          );
      });
    },
    workflowTags: function() {
      ref.current.workflows
        .selectAll(`.${styles.workflowTag}`)
        .classed(styles.published, published);
    }
  };

  const create = {
    workflows: function() {
      newWorkflows = ref.current.workflows
        .enter()
        .append('g')
        .classed(styles.workflow, true)
        .attr(
          'transform',
          (d: Workflow, idx: number) => `translate(0,${idx * 100})`
        );
    },
    nodes: function() {
      newNodes = nodes
        .enter()
        .append('g')
        .attr('id', (d: Node) => `node_${d.id}`)
        .classed(styles.node, true)
        .attr(
          'transform',
          (d: Node, idx: number) => `translate(${xScale(idx.toString())}, 0)`
        )
        .on('mouseenter', (d: Node) => events.nodeHighlight(d, true))
        .on('mouseleave', (d: Node) => events.nodeHighlight(d, false))
        .each(function(d: Node) {
          // @ts-ignore
          select(this)
            .append('g')
            .attr('transform', 'translate(0, 0)')
            .html(
              ReactDOMServer.renderToString(
                <VersionNode
                  type={d.type || TYPES.DEFAULT}
                  width={ref.current.nodeWidth}
                  height={DEFAULT_NODE_HEIGHT * ref.current.nodeSizeRatio}
                  status={d.status}
                />
              )
            );
        });
      newNodes
        .append('text')
        .classed('nodeText', true)
        .classed(styles.nodeText, true)
        .attr(
          'x',
          (d: Node) =>
            getNodeTextPadding(d.type || TYPES.DEFAULT) *
            ref.current.nodeSizeRatio
        )
        .attr('y', (DEFAULT_NODE_HEIGHT * ref.current.nodeSizeRatio) / 2)
        .attr('dy', 0)
        .style('font-size', `${fontSize}px`)
        .text((d: Node) => d.name)
        .call(wrap, 47 * ref.current.nodeSizeRatio)
        .call(centerText, fontSize);
    },
    workflowTags: function() {
      workflowsTag = newWorkflows
        .append('g')
        .attr('transform', `translate(0, 50)`)
        .classed(styles.published, published)
        .classed(styles.workflowTag, true);
      workflowsTag
        .append('text')
        .classed(styles.workflowTagText, true)
        .attr('x', marginWorkflow - 10)
        .attr('y', nodeHeight / 2)
        .attr('dy', 0)
        .style('font-size', `${fontSize}px`)
        .text((d: Workflow) => d.name)
        .call(centerText, fontSize);
      workflowsTag
        .append('line')
        .attr('x1', marginWorkflow)
        .attr('x2', marginWorkflow)
        .attr('y1', 2.5)
        .attr('y2', -2.5 + nodeHeight)
        .attr('stroke-width', STROKE_WIDTH);
      workflowsTag
        .append('line')
        .attr('x1', marginWorkflow)
        .attr('x2', xScale('0'))
        .attr('y1', (DEFAULT_NODE_HEIGHT * ref.current.nodeSizeRatio) / 2)
        .attr('y2', (DEFAULT_NODE_HEIGHT * ref.current.nodeSizeRatio) / 2)
        .attr('stroke-width', STROKE_WIDTH)
        .attr('stroke-dasharray', '3, 3');
      workflowsTag
        .append('path')
        .attr('d', 'M 0 0 l -7 -3 M 0 0 l -7 3')
        .attr(
          'transform',
          `translate(${xScale('0')}, ${(DEFAULT_NODE_HEIGHT *
            ref.current.nodeSizeRatio) /
            2})`
        )
        .attr('stroke-width', STROKE_WIDTH);
    },
    edges: function() {
      newEdges = edges
        .enter()
        .append('g')
        .attr('class', (d: Edge) => styles[d.status || NodeStatus.STARTED])
        .classed(styles.edge, true);
      newEdges
        .append('path')
        .classed(styles.edgeLine, true)
        .attr('d', 'M 0 0 m -7 -3 l +7 +3 l -7 3')
        .attr(
          'transform',
          (d: Edge) =>
            `translate(${xScale(
              nodeIdToIndex[d.toNode]
            )}, ${(DEFAULT_NODE_HEIGHT * ref.current.nodeSizeRatio) / 2})`
        )
        .attr('stroke-width', STROKE_WIDTH);
      newEdges
        .append('line')
        .classed('edgeLine', true)
        .classed(styles.edgeLine, true)
        .attr(
          'x1',
          // @ts-ignore
          (d: Edge) =>
            xScale(nodeIdToIndex[d.fromNode]) + ref.current.nodeWidth - 2
        )
        .attr('x2', (d: Edge) => xScale(nodeIdToIndex[d.toNode]))
        .attr('y1', (DEFAULT_NODE_HEIGHT * ref.current.nodeSizeRatio) / 2)
        .attr('y2', (DEFAULT_NODE_HEIGHT * ref.current.nodeSizeRatio) / 2)
        .attr('stroke-dasharray', '3, 3')
        .attr('stroke-width', STROKE_WIDTH);
      newEdges
        .append('line')
        .classed(styles.lineContainer, true)
        .attr(
          'x1',
          // @ts-ignore
          (d: Edge) =>
            xScale(nodeIdToIndex[d.fromNode]) + ref.current.nodeWidth - 2
        )
        .attr('x2', (d: Edge) => xScale(nodeIdToIndex[d.toNode]))
        .attr('y1', (DEFAULT_NODE_HEIGHT * ref.current.nodeSizeRatio) / 2)
        .attr('y2', (DEFAULT_NODE_HEIGHT * ref.current.nodeSizeRatio) / 2)
        .on('mouseenter', function() {
          // @ts-ignore
          events.edgeHighlight(this, true);
        })
        .on('mouseleave', function() {
          // @ts-ignore
          events.edgeHighlight(this, false);
        });
    },
    bus: function() {
      busG = g
        .append('g')
        .classed(styles.busG, true)
        .attr('transform', `translate(0, ${80 + data.length * 100})`);

      busG
        .append('text')
        .classed(styles.workflowTagText, true)
        .attr('x', 20)
        .attr('y', 0)
        .attr('dy', 0)
        .style('font-size', `${fontSize}px`)
        .text('BUS')
        .call(centerText, fontSize);
      busG
        .append('line')
        .classed(styles.busLine, true)
        .attr('x1', 30)
        .attr('x2', width - margin.right - outerPadding)
        .attr('y1', 0)
        .attr('y2', 0)
        .attr('stroke-width', STROKE_WIDTH);
      busG
        .append('path')
        .classed(styles.busLineArrow, true)
        .attr('d', 'M 0 0 m -7 -3 l +7 +3 l -7 3')
        .attr(
          'transform',
          `translate(${width - margin.right - outerPadding}, 0)`
        )
        .attr('stroke-width', STROKE_WIDTH);
    }
  };

  const events = {
    nodeHighlight: function(d: Node, enter: boolean = true) {
      const node = select(`#node_${d.id}`);
      const strokeWidth = enter ? 2 : 1,
        filter = enter ? 'url(#glow)' : 'none',
        addHoveredClass = enter;

      node.select('.int').attr('stroke-width', strokeWidth);
      node.select('.ext').attr('stroke-width', strokeWidth);
      node.select('.ext').style('filter', filter);
      node.select('.nodeText').classed(styles.hovered, addHoveredClass);
    },
    edgeHighlight: function(node: any, enter: boolean = true) {
      const strokeDashArray = enter ? '10, 0.01' : '3, 3',
        strokeWidth = enter ? 1 : 0.7;

      // @ts-ignore
      const lines = select(node.parentNode).selectAll(`.${styles.edgeLine}`);
      lines.attr('stroke-dasharray', strokeDashArray);
      lines.attr('stroke-width', strokeWidth);
    }
  };

  return (
    <div className={styles.container} ref={container}>
      <svg width={width} height={height} ref={svg} className={styles.wrapper} />
    </div>
  );
}

export default VersionStatusViewer;
