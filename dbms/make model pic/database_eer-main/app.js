// ==========================================================
// Mermaid erDiagram Parser
// ==========================================================
function parseMermaidER(code) {
    const entities = {};
    const relationships = [];
    const lines = code.split('\n');
    let currentEntity = null;
    let insideBlock = false;

    for (let i = 0; i < lines.length; i++) {
        let line = lines[i].trim();
        if (!line || line.startsWith('%%') || line === 'erDiagram') continue;

        if (!insideBlock) {
            const blockStart = line.match(/^(\w+)\s*\{/);
            if (blockStart) {
                currentEntity = blockStart[1];
                if (!entities[currentEntity]) {
                    entities[currentEntity] = { name: currentEntity, attributes: [] };
                }
                insideBlock = true;
                continue;
            }
        }

        if (insideBlock && line === '}') {
            insideBlock = false;
            currentEntity = null;
            continue;
        }

        if (insideBlock && currentEntity) {
            const attrMatch = line.match(/^(\S+)\s+(\S+)(?:\s+(PK|FK|"[^"]*"))?(?:\s+(PK|FK|"[^"]*"))?(?:\s+("[^"]*"))?/);
            if (attrMatch) {
                const attrType = attrMatch[1];
                const attrName = attrMatch[2];
                let isPK = false, isFK = false, comment = '';
                for (let j = 3; j <= 5; j++) {
                    if (!attrMatch[j]) continue;
                    if (attrMatch[j] === 'PK') isPK = true;
                    else if (attrMatch[j] === 'FK') isFK = true;
                    else if (attrMatch[j].startsWith('"')) comment = attrMatch[j].replace(/"/g, '');
                }
                entities[currentEntity].attributes.push({ name: attrName, type: attrType, isPK, isFK, comment });
            }
            continue;
        }

        const relMatch = line.match(/^(\w+)\s+(\|[|o]|[{}][|o]|\|\|)(--)(\|[|o]|[{}][|o]|\|\||[|o][{}]|[|o]\|)\s+(\w+)\s*:\s*"?([^"]*)"?$/);
        if (relMatch) {
            const src = relMatch[1];
            const leftCard = relMatch[2];
            const rightCard = relMatch[4];
            const tgt = relMatch[5];
            const label = relMatch[6].trim();
            if (!entities[src]) entities[src] = { name: src, attributes: [] };
            if (!entities[tgt]) entities[tgt] = { name: tgt, attributes: [] };

            const parseCard = (token) => {
                let many = false, mandatory = false;
                if (token.includes('{') || token.includes('}')) many = true;
                if (token.includes('||') || (token.includes('|') && !token.includes('o'))) mandatory = true;
                if (token === '||') { many = false; mandatory = true; }
                if (token === 'o|' || token === '|o') { many = false; mandatory = false; }
                if (token === 'o{' || token === '}o') { many = true; mandatory = false; }
                if (token === '|{' || token === '}|') { many = true; mandatory = true; }
                return { many, mandatory };
            };
            const left = parseCard(leftCard);
            const right = parseCard(rightCard);

            relationships.push({
                source: src, target: tgt, label,
                srcCard: left.many ? 'N' : '1',
                tgtCard: right.many ? 'N' : '1',
                srcMandatory: left.mandatory,
                tgtMandatory: right.mandatory,
                isSelfRef: src === tgt
            });
        }
    }

    const isSpecializationLabel = (label) => {
        const normalized = (label || '').toLowerCase();
        return normalized === 'is a' || normalized.startsWith('is a ');
    };

    const specializations = {};
    const normalRelationships = [];

    relationships.forEach(rel => {
        if (isSpecializationLabel(rel.label)) {
            const parent = rel.source;
            if (!specializations[parent]) {
                specializations[parent] = {
                    parent,
                    children: [],
                    constraint: '',
                    total: true
                };
            }

            const dMatch = rel.label.match(/\((d|o)\)/i) || rel.label.match(/\b(d|o)\b/i);
            if (dMatch) {
                specializations[parent].constraint = dMatch[1].toLowerCase();
            }
            specializations[parent].children.push(rel.target);
            specializations[parent].total = specializations[parent].total && rel.srcMandatory;
            return;
        }

        normalRelationships.push(rel);
    });

    // Detect weak entities: has attributes but none is PK
    Object.values(entities).forEach(ent => {
        ent.isWeak = ent.attributes.length > 0 && !ent.attributes.some(a => a.isPK);
    });

    // Detect identifying relationships: connected to at least one weak entity
    normalRelationships.forEach(rel => {
        rel.isIdentifying = !!(entities[rel.source]?.isWeak || entities[rel.target]?.isWeak);
    });

    return {
        entities: Object.values(entities),
        relationships: normalRelationships,
        specializations: Object.values(specializations)
    };
}

// ==========================================================
// D3 Renderer
// ==========================================================
let simulation = null;
let currentZoom = null;

function renderDiagram(parsed) {
    const svg = d3.select('#svgCanvas');
    svg.selectAll('g.root').remove();
    if (simulation) simulation.stop();

    const { entities, relationships, specializations = [] } = parsed;
    const nodesData = [];
    const linksData = [];
    const nodeMap = {};
    let totalAttrs = 0;

    // Entity nodes
    entities.forEach(ent => {
        const node = { id: ent.name, type: 'entity', label: ent.name, isWeak: ent.isWeak };
        nodesData.push(node);
        nodeMap[ent.name] = node;

        ent.attributes.forEach(attr => {
            const attrId = `attr__${ent.name}__${attr.name}`;
            const attrType = attr.isPK ? 'pk' : (attr.isFK ? 'fk' : 'attr');
            const attrNode = { id: attrId, type: attrType, label: attr.name, parent: ent.name };
            nodesData.push(attrNode);
            nodeMap[attrId] = attrNode;
            linksData.push({ source: ent.name, target: attrId, linkType: 'attr' });
            totalAttrs++;
        });
    });

    // Relationship nodes + links
    relationships.forEach((rel, idx) => {
        const relId = `rel__${idx}__${rel.label}`;
        const relNode = { id: relId, type: 'rel', label: rel.label, isIdentifying: rel.isIdentifying };
        nodesData.push(relNode);
        nodeMap[relId] = relNode;

        linksData.push({ source: rel.source, target: relId, linkType: 'rel', mandatory: rel.srcMandatory, cardLabel: rel.srcCard });

        if (rel.isSelfRef) {
            const dummyId = `dummy__${idx}`;
            nodesData.push({ id: dummyId, type: 'dummy', label: '' });
            linksData.push({ source: relId, target: dummyId, linkType: 'rel', mandatory: false, cardLabel: '' });
            linksData.push({ source: dummyId, target: rel.target, linkType: 'rel', mandatory: rel.tgtMandatory, cardLabel: rel.tgtCard });
        } else {
            linksData.push({ source: relId, target: rel.target, linkType: 'rel', mandatory: rel.tgtMandatory, cardLabel: rel.tgtCard });
        }
    });

    // Specialization (ISA) nodes + links
    specializations.forEach((spec, idx) => {
        const specId = `spec__${idx}__${spec.parent}`;
        const specNode = {
            id: specId,
            type: 'spec',
            label: spec.constraint || 'd',
            total: spec.total
        };
        nodesData.push(specNode);
        nodeMap[specId] = specNode;

        linksData.push({
            source: spec.parent,
            target: specId,
            linkType: 'spec',
            mandatory: spec.total,
            cardLabel: ''
        });

        spec.children.forEach(child => {
            linksData.push({
                source: specId,
                target: child,
                linkType: 'spec',
                mandatory: true,
                cardLabel: ''
            });
        });
    });

    const totalRelCount = relationships.length + specializations.reduce((sum, spec) => sum + spec.children.length, 0);

    document.getElementById('statEntities').textContent = entities.length;
    document.getElementById('statRelations').textContent = totalRelCount;
    document.getElementById('statAttrs').textContent = totalAttrs;

    const container = document.getElementById('canvasArea');
    const width = container.clientWidth;
    const height = container.clientHeight;

    const g = svg.append('g').attr('class', 'root');

    const zoom = d3.zoom()
        .scaleExtent([0.03, 4])
        .on('zoom', (event) => g.attr('transform', event.transform));
    svg.call(zoom);
    currentZoom = zoom;

    // Force simulation
    simulation = d3.forceSimulation(nodesData)
        .force('link', d3.forceLink(linksData).id(d => d.id).distance(d => {
            if (d.linkType === 'attr') return 38;
            if (d.linkType === 'spec') return 55;
            return 70;
        }))
        .force('charge', d3.forceManyBody().strength(d => {
            if (d.type === 'attr' || d.type === 'pk' || d.type === 'fk') return -15;
            if (d.type === 'rel') return -300;
            if (d.type === 'spec') return -170;
            if (d.type === 'dummy') return -20;
            return -600;
        }))
        .force('collide', d3.forceCollide().radius(d => {
            if (d.type === 'attr' || d.type === 'pk' || d.type === 'fk') return 22;
            if (d.type === 'dummy') return 5;
            if (d.type === 'spec') return 20;
            return 55;
        }))
        .force('center', d3.forceCenter(width / 2, height / 2))
        .alphaDecay(0.02);

    // ── Draw links ──────────────────────────────────────────────
    const linkGroup = g.append('g').attr('class', 'links').selectAll('g')
        .data(linksData).enter().append('g');

    // Attribute links (thin gray)
    linkGroup.filter(d => d.linkType === 'attr').append('line')
        .attr('stroke', '#c8d2e6')
        .attr('stroke-width', 1.2)
        .attr('opacity', 0.85);

    // Total participation: double-line (thick base + light gap)
    linkGroup.filter(d => d.linkType === 'rel' && d.mandatory).each(function() {
        d3.select(this).append('line').attr('class', 'link-bg')
            .attr('stroke', '#374151').attr('stroke-width', 5).attr('opacity', 0.9);
        d3.select(this).append('line').attr('class', 'link-fg')
            .attr('stroke', '#edf1f8').attr('stroke-width', 2.2);
    });

    // Partial participation: single line
    linkGroup.filter(d => d.linkType === 'rel' && !d.mandatory).append('line')
        .attr('stroke', '#9ca3af').attr('stroke-width', 1.8).attr('opacity', 0.9);

    // Specialization links
    linkGroup.filter(d => d.linkType === 'spec' && d.mandatory).append('line')
        .attr('stroke', '#7c3aed').attr('stroke-width', 3.2).attr('opacity', 0.95);

    linkGroup.filter(d => d.linkType === 'spec' && !d.mandatory).append('line')
        .attr('stroke', '#a78bfa').attr('stroke-width', 2.1).attr('opacity', 0.9);

    // Cardinality labels
    const cardLabels = linkGroup.filter(d => d.cardLabel).append('text')
        .attr('class', 'card-label')
        .text(d => d.cardLabel);

    // ── Draw nodes ──────────────────────────────────────────────
    const nodeGroup = g.append('g').attr('class', 'nodes').selectAll('g')
        .data(nodesData).enter().append('g')
        .call(d3.drag()
            .on('start', (event, d) => {
                if (!event.active) simulation.alphaTarget(0.3).restart();
                d.fx = d.x; d.fy = d.y;
            })
            .on('drag', (event, d) => { d.fx = event.x; d.fy = event.y; })
            .on('end', (event, d) => {
                if (!event.active) simulation.alphaTarget(0);
            })
        );

    // ── Entity: rectangle (strong) or double-rectangle (weak)
    nodeGroup.filter(d => d.type === 'entity').each(function(d) {
        const w = 132, h = 44;
        const sel = d3.select(this);

        sel.append('rect')
            .attr('width', w).attr('height', h)
            .attr('x', -w/2).attr('y', -h/2).attr('rx', 4)
            .attr('fill', d.isWeak ? '#ede9fe' : '#dbeafe')
            .attr('stroke', d.isWeak ? '#7c3aed' : '#2563eb')
            .attr('stroke-width', 2)
            .attr('filter', 'url(#nodeShadow)');

        if (d.isWeak) {
            const inset = 5;
            sel.append('rect')
                .attr('width', w - inset*2).attr('height', h - inset*2)
                .attr('x', -w/2 + inset).attr('y', -h/2 + inset).attr('rx', 2)
                .attr('fill', 'none')
                .attr('stroke', '#7c3aed')
                .attr('stroke-width', 1.5);
        }
    });

    // ── Relationship: diamond (regular) or double-diamond (identifying)
    nodeGroup.filter(d => d.type === 'rel').each(function(d) {
        const sel = d3.select(this);

        sel.append('polygon')
            .attr('points', '0,-35 74,0 0,35 -74,0')
            .attr('fill', '#fef3c7')
            .attr('stroke', '#b45309')
            .attr('stroke-width', 2)
            .attr('filter', 'url(#nodeShadow)');

        if (d.isIdentifying) {
            sel.append('polygon')
                .attr('points', '0,-26 57,0 0,26 -57,0')
                .attr('fill', 'none')
                .attr('stroke', '#b45309')
                .attr('stroke-width', 1.5);
        }
    });

    // ── Attribute ellipses
    nodeGroup.filter(d => d.type === 'attr').append('ellipse')
        .attr('rx', 46).attr('ry', 16)
        .attr('fill', '#f8fafc').attr('stroke', '#94a3b8').attr('stroke-width', 1.5);

    // ── Specialization constraint circle (d/o)
    nodeGroup.filter(d => d.type === 'spec').append('circle')
        .attr('r', 16)
        .attr('fill', '#ede9fe')
        .attr('stroke', '#7c3aed')
        .attr('stroke-width', 2);

    nodeGroup.filter(d => d.type === 'pk').append('ellipse')
        .attr('rx', 48).attr('ry', 17)
        .attr('fill', '#eff6ff').attr('stroke', '#3b82f6').attr('stroke-width', 2);

    nodeGroup.filter(d => d.type === 'fk').append('ellipse')
        .attr('rx', 46).attr('ry', 16)
        .attr('fill', '#f0fdf4').attr('stroke', '#16a34a').attr('stroke-width', 1.5)
        .attr('stroke-dasharray', '5,3');

    // ── Labels
    nodeGroup.filter(d => d.type === 'entity').append('text')
        .attr('class', d => d.isWeak ? 'node-label weak' : 'node-label')
        .text(d => d.label);

    nodeGroup.filter(d => d.type === 'rel').append('text')
        .attr('class', 'rel-label').text(d => d.label);

    nodeGroup.filter(d => d.type === 'attr').append('text')
        .attr('class', 'attr-label').text(d => d.label);

    nodeGroup.filter(d => d.type === 'pk').append('text')
        .attr('class', 'attr-label pk-label')
        .style('text-decoration', 'underline')
        .text(d => d.label);

    nodeGroup.filter(d => d.type === 'fk').append('text')
        .attr('class', 'attr-label fk-label')
        .text(d => d.label);

    nodeGroup.filter(d => d.type === 'spec').append('text')
        .attr('class', 'spec-label')
        .text(d => d.label || 'd');

    // ── Hover interactions ───────────────────────────────────────
    const adj = {};
    linksData.forEach(l => {
        const s = typeof l.source === 'object' ? l.source.id : l.source;
        const t = typeof l.target === 'object' ? l.target.id : l.target;
        adj[`${s},${t}`] = true;
        adj[`${t},${s}`] = true;
    });
    function isConn(a, b) { return adj[`${a.id},${b.id}`] || a.id === b.id; }

    nodeGroup.filter(d => d.type === 'entity' || d.type === 'rel')
        .style('cursor', 'pointer')
        .on('mouseover', function(event, d) {
            d3.select(this).select('rect, polygon').attr('filter', 'url(#nodeHover)');
            nodeGroup.style('opacity', o => isConn(d, o) ? 1 : 0.12);
            linkGroup.style('opacity', l => {
                const s = typeof l.source === 'object' ? l.source.id : l.source;
                const t = typeof l.target === 'object' ? l.target.id : l.target;
                return (s === d.id || t === d.id) ? 1 : 0.06;
            });
        })
        .on('mouseout', function() {
            d3.select(this).select('rect, polygon').attr('filter', 'url(#nodeShadow)');
            nodeGroup.style('opacity', 1);
            linkGroup.style('opacity', 1);
        });

    // ── Tick ──────────────────────────────────────────────────────
    simulation.on('tick', () => {
        linkGroup.selectAll('line')
            .attr('x1', d => d.source.x).attr('y1', d => d.source.y)
            .attr('x2', d => d.target.x).attr('y2', d => d.target.y);

        cardLabels
            .attr('x', d => d.source.x * 0.6 + d.target.x * 0.4)
            .attr('y', d => d.source.y * 0.6 + d.target.y * 0.4 - 11);

        nodeGroup.attr('transform', d => `translate(${d.x},${d.y})`);
    });

    setTimeout(() => fitToScreen(svg, zoom, width, height), 1200);
}

function fitToScreen(svg, zoom, width, height) {
    if (!svg || !zoom) return;
    const g = svg.select('g.root');
    if (g.empty()) return;
    const bbox = g.node().getBBox();
    if (!bbox || bbox.width === 0 || bbox.height === 0) return;
    const pad = 60;
    const scale = Math.min(
        (width - pad * 2) / bbox.width,
        (height - pad * 2) / bbox.height,
        1
    );
    const tx = width / 2 - scale * (bbox.x + bbox.width / 2);
    const ty = height / 2 - scale * (bbox.y + bbox.height / 2);
    svg.transition().duration(800).call(
        zoom.transform,
        d3.zoomIdentity.translate(tx, ty).scale(scale)
    );
}

// ==========================================================
// UI Wiring
// ==========================================================
const textarea = document.getElementById('mermaidInput');
const statusBar = document.getElementById('statusBar');
const statusText = document.getElementById('statusText');

function setStatus(type, msg) {
    statusBar.className = 'editor-status ' + type;
    statusText.textContent = msg;
}

function doRender() {
    const code = textarea.value.trim();
    if (!code) { setStatus('err', '請輸入 Mermaid erDiagram 語法'); return; }
    try {
        const parsed = parseMermaidER(code);
        if (parsed.entities.length === 0) { setStatus('err', '未偵測到任何實體，請確認語法正確'); return; }
        renderDiagram(parsed);
        const weakCount = parsed.entities.filter(e => e.isWeak).length;
        const specCount = (parsed.specializations || []).reduce((sum, spec) => sum + spec.children.length, 0);
        setStatus('ok', `成功渲染 ${parsed.entities.length} 個實體（${weakCount} 弱）、${parsed.relationships.length + specCount} 條關係`);
    } catch(e) {
        setStatus('err', '解析錯誤：' + e.message);
        console.error(e);
    }
}

document.getElementById('btnRender').addEventListener('click', doRender);

textarea.addEventListener('keydown', (e) => {
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') { e.preventDefault(); doRender(); }
    if (e.key === 'Tab') {
        e.preventDefault();
        const start = textarea.selectionStart, end = textarea.selectionEnd;
        textarea.value = textarea.value.substring(0, start) + '    ' + textarea.value.substring(end);
        textarea.selectionStart = textarea.selectionEnd = start + 4;
    }
});

document.getElementById('btnClear').addEventListener('click', () => {
    textarea.value = '';
    d3.select('#svgCanvas').selectAll('g.root').remove();
    if (simulation) simulation.stop();
    document.getElementById('statEntities').textContent = '0';
    document.getElementById('statRelations').textContent = '0';
    document.getElementById('statAttrs').textContent = '0';
    setStatus('idle', '已清空');
});

document.getElementById('toggleEditor').addEventListener('click', () => {
    const panel = document.getElementById('editorPanel');
    panel.classList.toggle('collapsed');
    document.getElementById('toggleEditor').textContent = panel.classList.contains('collapsed') ? '▶' : '◀';
});

document.getElementById('btnFit').addEventListener('click', () => {
    const container = document.getElementById('canvasArea');
    if (currentZoom) fitToScreen(d3.select('#svgCanvas'), currentZoom, container.clientWidth, container.clientHeight);
});

document.getElementById('btnUnpin').addEventListener('click', () => {
    if (!simulation) return;
    simulation.nodes().forEach(d => { d.fx = null; d.fy = null; });
    simulation.alpha(0.8).restart();
});

document.getElementById('btnZoomIn').addEventListener('click', () => {
    if (currentZoom) d3.select('#svgCanvas').transition().duration(300).call(currentZoom.scaleBy, 1.5);
});
document.getElementById('btnZoomOut').addEventListener('click', () => {
    if (currentZoom) d3.select('#svgCanvas').transition().duration(300).call(currentZoom.scaleBy, 0.67);
});

// Load schema and auto-render
fetch('schema.mermaid')
    .then(r => r.text())
    .then(t => { textarea.value = t; doRender(); })
    .catch(err => console.error('Failed to load schema.mermaid:', err));
