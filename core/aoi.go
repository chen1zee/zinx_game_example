package core

import (
	"fmt"
)

const (
	AoiMinX  int = 85
	AoiMaxX  int = 410
	AoiCntsX int = 10
	AoiMinY  int = 10
	AoiMaxY  int = 400
	AoiCntsY int = 20
)

/*
	AOI 管理模块
*/
type AOIManager struct {
	MinX  int           // 区域左边界坐标
	MaxX  int           // 区域有边界坐标
	CntsX int           // x 方向格子数量
	MinY  int           // 区域上边界坐标
	MaxY  int           // 区域下边界坐标
	CntsY int           // y方向格子数量
	grids map[int]*Grid // 当前区域中都有哪些格子， key=格子ID, value=格子对象
}

// 得到每个格子在 X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子在Y轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印信息方法
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager: \n minX: %d, maxX: %d, cntX: %d, minY: %d, maxY: %d, cntsY: %d \n Grids in AOI Manager: \n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 根据格子gID 得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGridsByGrid(gID int) (grids []*Grid) {
	tempGIDs := make([]int, 0) // 用来存储 base 一行
	// 判断gID 是否存在
	if _, ok := m.grids[gID]; !ok {
		return
	}
	// 将当前gid添加到 九宫格中
	tempGIDs = append(tempGIDs, gID)
	// 根据 gid得到当前格子所在 X轴编号
	idx := gID % m.CntsX
	// 判断当前idx 左边是否还有格子
	if idx > 0 {
		tempGIDs = append(tempGIDs, gID-1)
	}
	// 判断当前idx右边是否还有格子
	if idx < m.CntsX-1 {
		tempGIDs = append(tempGIDs, gID+1)
	}
	// 一行的所有 格子， 是否有 上下格子的 情况都是一样的，
	// 所以 判断某个格子的 上下情况， 即可求得 一行的 上下情况
	idy := tempGIDs[0] / m.CntsX
	up := false // 上下是否有 格子 flag
	down := false
	if idy > 0 {
		up = true
	} // 上方有
	if idy < m.CntsY-1 {
		down = true
	} // 下方
	for _, gID := range tempGIDs {
		grids = append(grids, m.grids[gID])
		if up {
			grids = append(grids, m.grids[(gID-m.CntsX)])
		}
		if down {
			grids = append(grids, m.grids[gID+m.CntsX])
		}
	}
	return
}

// 通过横纵坐标获取对应的 格子ID
func (m *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.gridWidth()
	gy := (int(y) - m.MinY) / m.gridLength()
	return gy*m.CntsX + gx
}

// 通过横纵坐标得到周边九宫格内的全部 PlayerIDs
func (m *AOIManager) GetPlayerIDsByPos(x, y float32) (playerIDs []int) {
	// 根据 横纵坐标得到当前坐标属于哪个格子ID
	gID := m.GetGIDByPos(x, y)
	// 根据格子ID 得到周边九宫格信息
	grids := m.GetSurroundGridsByGrid(gID)
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
		fmt.Printf("===> grid ID: %d, playerIds : %v ===", grid.GID, grid.GetPlayerIDs())
	}
	return
}

// 通过GID 获取当前格子的全部playerID
func (m *AOIManager) GetPlayerIdsByGrid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

// 移除一个格子中的PlayerID
func (m *AOIManager) RemovePlayerIdFromGrid(playerID, gID int) {
	m.grids[gID].Remove(playerID)
}

// 添加一个 PlayerID 到一个格子中
func (m *AOIManager) AddPlayerIdToGrid(playerID, gID int) {
	m.grids[gID].Add(playerID)
}

// 通过横纵坐标添加一个Player到一个格子中
func (m *AOIManager) AddToGridByPos(playerID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.grids[gID]
	grid.Add(playerID)
}

// 通过横纵坐标把一个Player从对应的格子中删除
func (m *AOIManager) RemoveFromGridByPos(playerID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.grids[gID]
	grid.Remove(playerID)
}

/*
	初始化一个AOI区域
*/
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}
	// 给AOI 初始化区域中所有的格子
	gridW := aoiMgr.gridWidth()
	gridL := aoiMgr.gridLength()
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 计算格子ID
			gid := y*cntsX + x
			// 初始化一个格子 放在 AOI map 中，key是当前格子ID
			aoiMgr.grids[gid] = NewGrid(
				gid,
				aoiMgr.MinX+x*gridW,
				aoiMgr.MinX+(x+1)*gridW,
				aoiMgr.MinY+y*gridL,
				aoiMgr.MinY+(y+1)*gridL)
		}
	}
	return aoiMgr
}
