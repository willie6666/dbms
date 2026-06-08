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
            const attrMatch = line.match(/^(\S+)\s+(\S+)/);
            if (attrMatch) {
                const attrType = attrMatch[1];
                let attrName = attrMatch[2];
                if(attrName.endsWith(',')) attrName = attrName.slice(0, -1);
                
                let isPK = false, isFK = false;
                if (line.includes('PK')) isPK = true;
                if (line.includes('FK')) isFK = true;
                
                entities[currentEntity].attributes.push({ name: attrName, type: attrType, isPK, isFK });
            }
            continue;
        }

        const relMatch = line.match(/^(\w+)\s+(\|[|o]|[{}][|o]|\|\|)(--)(\|[|o]|[{}][|o]|\|\||[|o][{}]|[|o]\|)\s+(\w+)/);
        if (relMatch) {
            const src = relMatch[1];
            const srcConn = relMatch[2];
            const tgtConn = relMatch[4];
            const tgt = relMatch[5];
            if (!entities[src]) entities[src] = { name: src, attributes: [] };
            if (!entities[tgt]) entities[tgt] = { name: tgt, attributes: [] };

            relationships.push({ source: src, srcConn, target: tgt, tgtConn });
        }
    }

    return { entities: Object.values(entities), relationships };
}

function renderSchema(parsed) {
    const textLayer = document.getElementById('textLayer');
    const arrowLayer = document.getElementById('arrowLayer');
    textLayer.innerHTML = '';
    arrowLayer.innerHTML = '<defs></defs>';
    
    if (parsed.entities.length === 0) return;

    // 1. Render HTML Grid Blocks
    parsed.entities.forEach(ent => {
        const row = document.createElement('div');
        row.className = 'schema-row';
        row.id = `table-${ent.name}`;
        
        const titleSpan = document.createElement('div');
        titleSpan.className = 'table-name';
        titleSpan.textContent = ent.name;
        
        const attrContainer = document.createElement('div');
        attrContainer.className = 'attr-container';
        
        ent.attributes.forEach(attr => {
            const attrWrapper = document.createElement('div');
            attrWrapper.className = 'attr-wrapper';
            attrWrapper.id = `attr-${ent.name}-${attr.name}`;
            
            const attrSpan = document.createElement('span');
            attrSpan.className = 'attr';
            if (attr.isPK) attrSpan.classList.add('pk');
            if (attr.isFK) attrSpan.classList.add('fk');
            attrSpan.textContent = attr.name;
            
            attrWrapper.appendChild(attrSpan);
            attrContainer.appendChild(attrWrapper);
        });
        
        row.appendChild(titleSpan);
        row.appendChild(attrContainer);
        textLayer.appendChild(row);
    });

    // Wait for DOM to render to calculate positions
    setTimeout(() => {
        drawArrows(parsed);
    }, 100);
}

function drawArrows(parsed) {
    const arrowLayer = document.getElementById('arrowLayer');
    const containerRect = document.getElementById('schemaContainer').getBoundingClientRect();
    const defs = arrowLayer.querySelector('defs');
    
    const connections = [];
    
    parsed.entities.forEach(ent => {
        ent.attributes.filter(a => a.isFK).forEach(fkAttr => {
            let targetEnt = null;
            let targetPk = null;
            
            // 0. Hardcoded explicit dictionary mapping for irregular FKs
            const hardcodedFkMap = {
                'developer_id': ['USERS', 'USER'],
                'buyer_id': ['USERS', 'USER'],
                'handled_by': ['USERS', 'USER'],
                'sender_id': ['USERS', 'USER'],
                'receiver_id': ['USERS', 'USER'],
                'blocker_id': ['USERS', 'USER'],
                'blocked_id': ['USERS', 'USER'],
                'parent_reply_id': ['REVIEW_REPLY', 'REVIEW_REPLIES'],
                'transaction_item_id': ['TRANSACTION_ITEM', 'TRANSACTION_ITEMS']
            };
            
            if (hardcodedFkMap[fkAttr.name]) {
                const targetNames = hardcodedFkMap[fkAttr.name];
                for (let tName of targetNames) {
                    const targetEntity = parsed.entities.find(e => e.name === tName);
                    if (targetEntity) {
                        const pkMatch = targetEntity.attributes.find(a => a.isPK);
                        if (pkMatch) {
                            targetEnt = targetEntity;
                            targetPk = pkMatch;
                            break;
                        }
                    }
                }
            }
            
            const relatedRels = parsed.relationships.filter(r => r.source === ent.name || r.target === ent.name);
            
            // 1. Exact match
            if (!targetEnt) {
                for (let rel of relatedRels) {
                    const otherName = rel.source === ent.name ? rel.target : rel.source;
                    const otherEnt = parsed.entities.find(e => e.name === otherName);
                    if (otherEnt) {
                        let pkMatch = otherEnt.attributes.find(a => a.isPK && a.name === fkAttr.name);
                        if (pkMatch) {
                            targetEnt = otherEnt; targetPk = pkMatch; break;
                        }
                    }
                }
            }

            // 2. Partial match
            if (!targetEnt) {
                for (let rel of relatedRels) {
                    const otherName = rel.source === ent.name ? rel.target : rel.source;
                    const otherEnt = parsed.entities.find(e => e.name === otherName);
                    if (otherEnt) {
                        let pkMatch = otherEnt.attributes.find(a => a.isPK && (fkAttr.name.includes(a.name.replace('_id','')) || a.name.includes(fkAttr.name.replace('_id',''))));
                        if (pkMatch) {
                            targetEnt = otherEnt; targetPk = pkMatch; break;
                        }
                    }
                }
            }

            // 3. Global match
            if (!targetEnt) {
                for (let otherEnt of parsed.entities) {
                    if (otherEnt.name === ent.name) continue;
                    const pkMatch = otherEnt.attributes.find(a => a.isPK && a.name === fkAttr.name);
                    if (pkMatch) {
                        targetEnt = otherEnt; targetPk = pkMatch; break;
                    }
                }
            }
            
            // 4. Fallback (Only if this entity is the CHILD side in exactly one parent relationship)
            if (!targetEnt) {
                const parentRels = parsed.relationships.filter(r => {
                    if (r.target === ent.name && r.tgtConn.includes('{')) return true;
                    if (r.source === ent.name && r.srcConn.includes('}')) return true;
                    return false;
                });
                
                if (parentRels.length === 1) {
                    const parentName = parentRels[0].source === ent.name ? parentRels[0].target : parentRels[0].source;
                    const parentEnt = parsed.entities.find(e => e.name === parentName);
                    if (parentEnt) {
                        let pkMatch = parentEnt.attributes.find(a => a.isPK);
                        if (pkMatch) {
                            targetEnt = parentEnt; targetPk = pkMatch;
                        }
                    }
                }
            }
            
            if (targetEnt && targetPk) {
                connections.push({
                    fkTable: ent.name,
                    fkAttr: fkAttr.name,
                    pkTable: targetEnt.name,
                    pkAttr: targetPk.name
                });
            }
        });
    });

    const colors = ['#2c3e50', '#8e44ad', '#2980b9', '#16a085', '#d35400', '#c0392b', '#7f8c8d'];
    let leftMarginCounter = 60;
    let rightMarginCounter = document.getElementById('textLayer').offsetWidth + 100;

    connections.forEach((conn, idx) => {
        const fkElem = document.getElementById(`attr-${conn.fkTable}-${conn.fkAttr}`);
        const pkElem = document.getElementById(`attr-${conn.pkTable}-${conn.pkAttr}`);
        
        if (!fkElem || !pkElem) return;

        const fkRect = fkElem.getBoundingClientRect();
        const pkRect = pkElem.getBoundingClientRect();
        const fkRowRect = fkElem.closest('.schema-row').getBoundingClientRect();
        const pkRowRect = pkElem.closest('.schema-row').getBoundingClientRect();
        
        const fkIndex = parsed.entities.findIndex(e => e.name === conn.fkTable);
        const pkIndex = parsed.entities.findIndex(e => e.name === conn.pkTable);
        
        const x1 = fkRect.left + fkRect.width / 2 - containerRect.left;
        const y1_bottom = fkRect.bottom - containerRect.top;
        
        const x2 = pkRect.left + pkRect.width / 2 - containerRect.left;
        const y2_bottom = pkRect.bottom - containerRect.top;
        
        const color = colors[idx % colors.length];
        const isSelf = fkIndex === pkIndex;
        
        const goLeft = x1 < containerRect.width / 2;
        let routeX;
        
        if (goLeft) {
            routeX = leftMarginCounter;
            leftMarginCounter -= 10;
        } else {
            routeX = rightMarginCounter;
            rightMarginCounter += 10;
        }

        let d = '';
        const yOffset = (idx % 6) * 5; // staggering vertical drops below tables

        if (isSelf) {
            // U-turn below the table
            const midY = fkRowRect.bottom - containerRect.top + 10 + yOffset;
            d = `M ${x1},${y1_bottom} L ${x1},${midY} L ${x2},${midY} L ${x2},${y2_bottom}`;
        } else {
            // Classic Chen layout: exit bottom, route below tables, enter bottom
            const midY_FK = fkRowRect.bottom - containerRect.top + 10 + yOffset;
            const midY_PK = pkRowRect.bottom - containerRect.top + 10 + yOffset;
            d = `M ${x1},${y1_bottom} L ${x1},${midY_FK} L ${routeX},${midY_FK} L ${routeX},${midY_PK} L ${x2},${midY_PK} L ${x2},${y2_bottom}`;
        }

        const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
        path.setAttribute("d", d);
        path.setAttribute("stroke", color);
        
        const markerId = `arrowhead-${idx}`;
        const marker = document.createElementNS("http://www.w3.org/2000/svg", "marker");
        marker.setAttribute("id", markerId);
        marker.setAttribute("markerWidth", "7");
        marker.setAttribute("markerHeight", "7");
        marker.setAttribute("refX", "7");
        marker.setAttribute("refY", "3.5");
        marker.setAttribute("orient", "auto");
        
        const poly = document.createElementNS("http://www.w3.org/2000/svg", "polygon");
        poly.setAttribute("points", "0,0 0,7 7,3.5");
        poly.setAttribute("fill", color);
        
        marker.appendChild(poly);
        defs.appendChild(marker);
        
        path.setAttribute("marker-end", `url(#${markerId})`);
        
        // Add hover effect
        path.addEventListener('mouseenter', () => { path.style.strokeWidth = '3'; });
        path.addEventListener('mouseleave', () => { path.style.strokeWidth = '1.5'; });
        
        arrowLayer.appendChild(path);
    });
}

document.getElementById('renderBtn').addEventListener('click', () => {
    const input = document.getElementById('mermaidInput').value;
    const parsed = parseMermaidER(input);
    renderSchema(parsed);
});

document.getElementById('clearBtn').addEventListener('click', () => {
    document.getElementById('mermaidInput').value = '';
    document.getElementById('textLayer').innerHTML = '';
    document.getElementById('arrowLayer').innerHTML = '<defs></defs>';
});

// Load example
fetch('../database_eer-main/schema.mermaid')
    .then(response => {
        if(response.ok) return response.text();
        return null;
    })
    .then(text => {
        if (text) {
            document.getElementById('mermaidInput').value = text;
            setTimeout(() => {
                document.getElementById('renderBtn').click();
            }, 300);
        }
    })
    .catch(() => {});
