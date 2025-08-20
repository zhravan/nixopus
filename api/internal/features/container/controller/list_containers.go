package controller

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/go-fuego/fuego"
	containertypes "github.com/raghavyuva/nixopus-api/internal/features/container/types"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

func (c *ContainerController) ListContainers(fuegoCtx fuego.ContextNoBody) (*shared_types.Response, error) {
	// normalize query params
	params := parseContainerListParams(fuegoCtx.Request())

	// Get pre-filtered summaries from Docker
	containers, err := c.dockerService.ListContainers(container.ListOptions{
		All:     true,
		Filters: buildDockerFilters(params),
	})
	if err != nil {
		c.logger.Log(logger.Error, err.Error(), "")
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}
	// Build summaries, then search/sort/paginate
	rows := summarizeContainers(containers)
	pageRows, totalCount := applySearchSortPaginate(rows, params)

	result := c.enrichContainers(pageRows, containers)

	return &shared_types.Response{
		Status:  "success",
		Message: "Containers fetched successfully",
		Data: map[string]interface{}{
			"containers":  result,
			"total_count": totalCount,
			"page":        params.page,
			"page_size":   params.pageSize,
			"sort_by":     params.sortBy,
			"sort_order":  params.sortOrder,
			"search":      params.search,
			"status":      params.status,
			"name":        params.name,
			"image":       params.image,
		},
	}, nil
}

type containerListParams struct {
	page      int
	pageSize  int
	search    string
	sortBy    string
	sortOrder string
	status    string
	name      string
	image     string
}

func parseContainerListParams(r *http.Request) containerListParams {
	q := r.URL.Query()
	pageStr := q.Get("page")
	pageSizeStr := q.Get("page_size")
	sortBy := strings.ToLower(strings.TrimSpace(q.Get("sort_by")))
	sortOrder := strings.ToLower(strings.TrimSpace(q.Get("sort_order")))

	if pageStr == "" {
		pageStr = "1"
	}
	if pageSizeStr == "" {
		pageSizeStr = "10"
	}
	if sortBy == "" {
		sortBy = "name"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 10
	}

	return containerListParams{
		page:      page,
		pageSize:  pageSize,
		search:    strings.TrimSpace(q.Get("search")),
		sortBy:    sortBy,
		sortOrder: sortOrder,
		status:    strings.TrimSpace(q.Get("status")),
		name:      strings.TrimSpace(q.Get("name")),
		image:     strings.TrimSpace(q.Get("image")),
	}
}

func buildDockerFilters(p containerListParams) filters.Args {
	f := filters.NewArgs()
	if p.status != "" {
		f.Add("status", p.status)
	}
	if p.name != "" {
		f.Add("name", p.name)
	}
	if p.image != "" {
		f.Add("ancestor", p.image)
	}
	return f
}

type listRow struct {
	ID      string
	Name    string
	Image   string
	Status  string
	State   string
	Created int64
	Labels  map[string]string
}

func summarizeContainers(summaries []container.Summary) []listRow {
	rows := make([]listRow, 0, len(summaries))
	for _, csum := range summaries {
		name := ""
		if len(csum.Names) > 0 {
			n := csum.Names[0]
			if len(n) > 1 {
				name = n[1:]
			} else {
				name = n
			}
		}
		rows = append(rows, listRow{
			ID:      csum.ID,
			Name:    name,
			Image:   csum.Image,
			Status:  csum.Status,
			State:   csum.State,
			Created: csum.Created,
			Labels:  csum.Labels,
		})
	}
	return rows
}

func applySearchSortPaginate(rows []listRow, p containerListParams) ([]listRow, int) {
	if p.search != "" {
		lower := strings.ToLower(p.search)
		filtered := make([]listRow, 0, len(rows))
		for _, r := range rows {
			if strings.Contains(strings.ToLower(r.Name), lower) ||
				strings.Contains(strings.ToLower(r.Image), lower) ||
				strings.Contains(strings.ToLower(r.Status), lower) {
				filtered = append(filtered, r)
			}
		}
		rows = filtered
	}

	sort.Slice(rows, func(i, j int) bool {
		var less bool
		switch p.sortBy {
		case "status":
			less = strings.ToLower(rows[i].Status) < strings.ToLower(rows[j].Status)
		case "name":
			less = strings.ToLower(rows[i].Name) < strings.ToLower(rows[j].Name)
		default:
			less = rows[i].Created < rows[j].Created
		}
		if p.sortOrder == "desc" {
			return !less
		}
		return less
	})

	totalCount := len(rows)
	start := (p.page - 1) * p.pageSize
	if start > totalCount {
		start = totalCount
	}
	end := start + p.pageSize
	if end > totalCount {
		end = totalCount
	}
	return rows[start:end], totalCount
}

func (c *ContainerController) enrichContainers(pageRows []listRow, summaries []container.Summary) []containertypes.Container {
	result := make([]containertypes.Container, 0, len(pageRows))
	for _, r := range pageRows {
		info, err := c.dockerService.GetContainerById(r.ID)
		if err != nil {
			c.logger.Log(logger.Error, "Error inspecting container", r.ID)
			continue
		}
		cd := containertypes.Container{
			ID:        r.ID,
			Name:      r.Name,
			Image:     r.Image,
			Status:    r.Status,
			State:     r.State,
			Created:   info.Created,
			Labels:    r.Labels,
			Command:   "",
			IPAddress: info.NetworkSettings.IPAddress,
			HostConfig: containertypes.HostConfig{
				Memory:     info.HostConfig.Memory,
				MemorySwap: info.HostConfig.MemorySwap,
				CPUShares:  info.HostConfig.CPUShares,
			},
		}
		if info.Config != nil && info.Config.Cmd != nil && len(info.Config.Cmd) > 0 {
			cd.Command = info.Config.Cmd[0]
		}
		for _, s := range summaries {
			if s.ID == r.ID {
				for _, p := range s.Ports {
					cd.Ports = append(cd.Ports, containertypes.Port{
						PrivatePort: int(p.PrivatePort),
						PublicPort:  int(p.PublicPort),
						Type:        p.Type,
					})
				}
				break
			}
		}
		for _, m := range info.Mounts {
			cd.Mounts = append(cd.Mounts, containertypes.Mount{
				Type:        string(m.Type),
				Source:      m.Source,
				Destination: m.Destination,
				Mode:        m.Mode,
			})
		}
		for name, network := range info.NetworkSettings.Networks {
			cd.Networks = append(cd.Networks, containertypes.Network{
				Name:       name,
				IPAddress:  network.IPAddress,
				Gateway:    network.Gateway,
				MacAddress: network.MacAddress,
				Aliases:    network.Aliases,
			})
		}
		result = append(result, cd)
	}
	return result
}
